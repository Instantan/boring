import * as sass from 'sass';

function global(): {
    location: {
        href: string
    },
    compileScssToCss: (str: string) => string,
    canonicalize: (str: string) => string,
    loadFile: (str: string) => { contents: string, syntax: "scss" | "css" | "indented" } | null
} {
    return eval('globalThis')
}

export function compileScssToCss(str: string) {
    return sass.compileString(str, {
        style: "compressed",
        importers: [
            {
                canonicalize(url) {
                    return new URL(global().canonicalize(url));
                },
                load(canonicalUrl) {
                    const urlStr = canonicalUrl.toString() || ""
                    const loadFile = global().loadFile
                    if (loadFile) {
                        return loadFile(urlStr)
                    }
                    return null
                }
            }
        ],
    }).css
}

// location.href polyfill
global().location = { href: "/" }

// registering functions
global().compileScssToCss = compileScssToCss