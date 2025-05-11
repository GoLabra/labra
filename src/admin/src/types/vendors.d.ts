declare module 'save-as' {
    function saveAs(blob: Blob, filename: string): void;
    export default saveAs;
}

declare module 'graphql-prettier' {
    function gqlPrettier(query: string): string;
    export default gqlPrettier;
}
