import gqlPrettier from 'graphql-prettier';

export class PostmanCollectionBuilder {

    private _postmanId: string | null = null;
    private _name: string | null = null;
    private _schema = 'https://schema.getpostman.com/json/collection/v2.1.0/collection.json';
    private _exporterId = '696969';
    private _items: PostmanItemBase[] = [];

    public addPostmanId = (postmanId: string) => {
        this._postmanId = postmanId;
        return this;
    }

    public addName = (name: string) => {
        this._name = name;
        return this;
    }

    public addItem = (item: PostmanItemBase) => {
        this._items.push(item);
        return this;
    }
    
    public addItems = (items: PostmanItemBase[]) => {
        this._items.push(...items);
        return this;
    }

    public build = () => {
        return {
            "info": {
                "_postman_id": this._postmanId,
                "name": this._name,
                "schema": this._schema,
                "_exporter_id": this._exporterId
            },
            item: this._items.map(i => i.build())
        }
    }
}


class PostmanItemBase {
    
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

export class PostmanDirectory extends PostmanItemBase {
    private _items: PostmanItem[] = [];

    public addItem = (item: PostmanItem) => {
        this._items.push(item);
        return this;
    }

    protected buildData = (): any => {
        return {
            name: this._name,
            item: this._items.map(i => i.build())
        }
    }
}


export class PostmanItem extends PostmanItemBase {

    private _auth = 'noauth';
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

        const url = new URL(this._url!);

        return {
            name: this._name,
            request: {
                auth: {
                    "type": this._auth
                },
                method: "POST",
                header: [],
                body:  {
                    "mode": "graphql",
                    "graphql": {
                        "query": gqlPrettier(this._query!),
                        "variables": JSON.stringify(this._variables, null, 2)
                    }
                },
                "url": {
                    "raw": this._url,
                    "protocol": "http",
                    "host": [
                        url.hostname
                    ],
                    "port": url.port,
                    "path": [
                        url.pathname
                    ]
                }
            }
        };
    }
}