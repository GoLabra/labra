
const systemKeywords = ['entity', 'ent', 'break' ,'default' ,'func' ,'interface' ,'select' ,'case' ,'defer' ,'go' ,'map' ,'struct' ,'chan' ,'else' ,'goto' ,'package' ,'switch' ,'const' ,'fallthrough' ,'if' ,'range' ,'type' ,'continue' ,'for' ,'import' ,'return' ,'var', 'int', 'string', 'bool', 'float', 'double', 'byte', 'rune', 'uint', 'int8', 'int16', 'int32', 'int64', 'uint8', 'uint16', 'uint32', 'uint64', 'uintptr', 'float32', 'float64', 'complex64', 'complex128', 'error', 'nil', 'true', 'false', 'iota', 'append', 'cap', 'close', 'complex', 'copy', 'delete', 'imag', 'len', 'make', 'new', 'panic', 'print', 'println', 'real', 'recover', 'order'];

export const SYSTEM_FIELDS = ['id', 'createdAt', 'updatedAt']; 
export const SYSTEM_EDGES= ['createdBy', 'updatedBy']; 
export const SYSTEM_CHILDREN = [ ...SYSTEM_FIELDS, ...SYSTEM_EDGES]

export const ENTITY_CHILDREN_SYSTEM_KEYWORDS = [ ...systemKeywords, ...SYSTEM_CHILDREN]
export const ENTITY_SYSTEM_KEYWORDS = [ ...systemKeywords]

export const BASE_PATH = process.env.BASE_PATH ?? ''

const runtimeEnv = await (async () => {
    if(!process.env.RUNTIME_ENV){
		return {};
	}

	if(typeof window == 'undefined'){
		return {}
	}

	const response = await fetch(`${BASE_PATH}/env.runtime.json`);
	if(!response.ok){
		return {}
	}
	
	const data = response.json();
	return data;
})();

export const GRAPHQL_API_URL = runtimeEnv.NEXT_PUBLIC_GRAPHQL_API_URL ?? process.env.NEXT_PUBLIC_GRAPHQL_API_URL ?? 'http://localhost:4000';

export const GRAPHQL_QUERY_API_URL = runtimeEnv.NEXT_PUBLIC_GRAPHQL_QUERY_API_URL ?? process.env.NEXT_PUBLIC_GRAPHQL_QUERY_API_URL ?? 'http://localhost:4000/query';
export const GRAPHQL_QUERY_PLAYGROUND_URL = runtimeEnv.NEXT_PUBLIC_GRAPHQL_QUERY_PLAYGROUND_URL ?? process.env.NEXT_PUBLIC_GRAPHQL_QUERY_PLAYGROUND_URL ?? 'http://localhost:4000/playground';

export const GRAPHQL_ADMIN_API_URL = runtimeEnv.NEXT_PUBLIC_GRAPHQL_ADMIN_API_URL ?? process.env.NEXT_PUBLIC_GRAPHQL_ADMIN_API_URL ?? 'http://localhost:4000/admin/query';
export const GRAPHQL_ADMIN_PLAYGROUND_URL = runtimeEnv.NEXT_PUBLIC_GRAPHQL_ADMIN_PLAYGROUND_URL ?? process.env.NEXT_PUBLIC_GRAPHQL_ADMIN_PLAYGROUND_URL ?? 'http://localhost:4000/admin/playground';

export const CENTRIFUGO_URL = runtimeEnv.NEXT_PUBLIC_CENTRIFUGO_URL ?? process.env.NEXT_PUBLIC_CENTRIFUGO_URL;

export const PRODUCT_NAME = runtimeEnv.NEXT_PUBLIC_BRAND_PRODUCT_NAME ?? process.env.NEXT_PUBLIC_BRAND_PRODUCT_NAME ?? 'Labra·GO';
export const COLOR = runtimeEnv.NEXT_PUBLIC_BRAND_COLOR ?? process.env.NEXT_PUBLIC_BRAND_COLOR ?? 'blue';

console.log(PRODUCT_NAME);
