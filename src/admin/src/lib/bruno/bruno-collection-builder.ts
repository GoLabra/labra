import gqlPrettier from 'graphql-prettier';

export class BrunoCollectionBuilder {

    private _postmanId: string | null = null;
    private _name: string | null = null;
    private _items: BrunoItemBase[] = [];

    public addPostmanId = (postmanId: string) => {
        this._postmanId = postmanId;
        return this;
    }

    public addName = (name: string) => {
        this._name = name;
        return this;
    }

    public addItem = (item: BrunoItemBase) => {
        this._items.push(item);
        return this;
    }

    public addItems = (items: BrunoItemBase[]) => {
        this._items.push(...items);
        return this;
    }

    public build = () => {
        return {
            "name": this._name,
            "version": "1",
            "environments": [],
            "brunoConfig": {
                "version": "1",
                "name": this._name,
                "type": "collection",
                "ignore": [
                    "node_modules",
                    ".git"
                ]
            },
            items: this._items.map(i => i.build())
        }
    }
}


class BrunoItemBase {

    protected _name: string | null = null;

    public addName = (name: string) => {
        this._name = name;
        return this;
    }

    protected buildData = (): any => {
        throw new Error("Method not implemented.");
    }

    public build = () => {
        return this.buildData();
    }
}

export class BrunoDirectory extends BrunoItemBase {
    private _items: BrunoItem[] = [];

    public addItem = (item: BrunoItem) => {
        this._items.push(item);
        return this;
    }

    protected buildData = (): any => {
        return {
            type: "folder",
            name: this._name,
            items: this._items.map(i => i.build())
        }
    }
}


export class BrunoItem extends BrunoItemBase {

    private _auth = 'none';
    private _url: string | null = null;
    private _query: string | null = null;
    private _variables: Record<string, any> | null = null;

    public addAuth = (auth: 'noauth' | 'basic' | 'bearer') => {
        this._auth = auth;
        return this;
    }

    public addUrl = (url: string) => {
        this._url = url;
        return this;
    }

    public addQuery = (query: string) => {
        this._query = query;
        return this;
    }

    public addVariables = (variables: Record<string, any>) => {
        this._variables = variables;
        return this;
    }

    public addQueryVariables = (query: { query: string, variables: Record<string, any> } | null) => {
        this._query = query?.query ?? null;
        this._variables = query?.variables ?? null;
        return this;
    }

    protected buildData = (): any => {

        return {
            "type": "graphql",
            "name": this._name,
            "seq": 1,
            "request": {
                "url": this._url,
                "method": "POST",
                "headers": [],
                "params": [],
                "body": {
                    "mode": "graphql",
                    "graphql": {
                        "query": gqlPrettier(this._query!),
                        "variables": JSON.stringify(this._variables, null, 2)
                    },
                    "formUrlEncoded": [],
                    "multipartForm": []
                },
                "script": {},
                "vars": {},
                "assertions": [],
                "tests": "",
                "docs": "",
                "auth": {
                    "mode": this._auth
                }
            }
        };
    }
}