// jest.config.js
/** @type {import('ts-jest/dist/types').InitialOptionsTsJest} */
module.exports = {
    // Use babel-jest for both JS and TS transformations
    // preset: "ts-jest", // removed, using babel-jest for TS as well
    testEnvironment: "jest-environment-jsdom",
    transform: {
        '^.+\\.[tj]sx?$': ['babel-jest', { configFile: './babel-jest.config.json' }],
        '^.+\\.css$': 'jest-scss-transform'
    },
    moduleFileExtensions: ['ts', 'tsx', 'js', 'jsx', 'json', 'node'],
    moduleNameMapper: {
        "^@/(.*)$": "<rootDir>/src/$1"
    },
    transformIgnorePatterns: [
        "node_modules/(?!(change-case|@mui/material/utils|ssr-window|dom7|@apollo|superjson)).*/"
    ]
}