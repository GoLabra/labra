import { Edge, EntityOwner } from "@/lib/apollo/graphql.entities";
import { useFullEntity } from "./use-entities";
import { useMemo } from "react";
import { gql, useQuery } from "@apollo/client";
import { ADMIN_CONTEXT } from "@/lib/apollo/apolloWrapper";
import { LGQuery } from "@/lib/apollo/builders/LabraGqlApiBuilder/LGQuery";
import { GplFilter } from "@/lib/apollo/builders/LabraGqlApiBuilder/types/types";
import { useLgQuery } from "./use-lg-query";
import { LGSelectInclude } from "@/lib/apollo/builders/LabraGqlApiBuilder/LGSelectInclude";

const GET_FILES_CONTENT = gql`query files($or: [FileWhereInput!]) {
	files(where: { or: $or }) {
		id
		content
	}
}`

export type LFile = {
	id: string;
	name: string;
	caption: string;
	mimeType: string;
	size: number;
}

interface UseEntityFilesParams {
	entityName: string;
	entryId: string | null | undefined;
	edge: Edge;
	loadContent?: (file: LFile) => boolean;
}
export const useEntityFiles = (props: UseEntityFilesParams) => {

	const edgeValueData = useRelationContentManagerStore<any[] | any>({
		entityName: props.entityName, 
		entryId: props.entryId, 
		edge: props.edge,
		fields: ['id', 'caption', 'name', 'mimeType', 'size'] 
	});

	const files = useMemo(() => {
		if(!edgeValueData.data){
			return [];
		}
		
		if(!Array.isArray(edgeValueData.data)){
			return [edgeValueData.data]
		}

		return edgeValueData.data ?? [];
	}, [edgeValueData.data]);

	const previewIds = useMemo(() => {
		return files.filter(i => props.loadContent?.(i) ?? true)
					.map(i => i.id);

	}, [files]);

	const filesContentData = useQuery<{ files: [{ id: string, content: string }] }>(GET_FILES_CONTENT, {
		variables: {
			or: previewIds
			.map(i => ({
				id: i
			})) ?? []
		}, 
		fetchPolicy: 'network-only',
		skip: previewIds.length == 0,
		context: ADMIN_CONTEXT
	});


	const filesWithContent = useMemo((): any[] => {
		
		const filesContent = filesContentData.data?.files ?? [];

		if(!filesContent.length){
			return files;
		}

		return files.map(file => {
			const contentFile = filesContent.find(f => f.id === file.id);
			return {
				...file,
				content: contentFile?.content
			};
		});
	}, [files, filesContentData.data])

	return filesWithContent;
}


interface UseRelationContentManagerStoreParams {
	entityName: string;
	entryId: string | null | undefined;
	edge: Edge;
	fields: 'iddisplay' | 'grid' | string[];
}
export const useRelationContentManagerStore = <T = any>(props: UseRelationContentManagerStoreParams) => {

	const rootEntity = useFullEntity({ entityName: props.entityName });
	const edgeEntity = useFullEntity({ entityName: props.edge.relatedEntity.name });


	const dataQuery = useMemo(() => {
		if(rootEntity?.loading ?? true){
			return null;
		}

		if(edgeEntity?.loading ?? true){
			return null;
		}

		if(!props.entryId){
			return null;
		}
		
		let query = LGQuery.from<any>(rootEntity!.name)
							.where(GplFilter.field('id', '', props.entryId));
							
		if(props.fields === 'iddisplay'){
			query = query.include(edgeEntity!.name, q => q.select('id', edgeEntity!.displayField!.name));
		} else if (Array.isArray(props.fields) && props.fields.every(item => typeof item === 'string')) {
			query = query.include(edgeEntity!.name, q => q.select(...props.fields));
		} else if(props.fields === 'grid'){
			
			query = query.include(props.edge.name, q => {
				// add fields
				let include = q.select(...edgeEntity!.fields.map(i => i.name));

				// add edges	
				include = edgeEntity!.edges.reduce((query: LGSelectInclude<any>, edge) => {
					return include.include(edge.name, q => q.select('id', edge.relatedEntity.displayField.name));
				}, include);

				return include;
			});			
		} 
		return query;	
	}, [rootEntity, edgeEntity, props.entryId, props.fields]);

	const apiType = rootEntity?.owner == EntityOwner.Admin ? 'admin' : 'user';
	const data = useLgQuery({
		apiType,
		query: [dataQuery],
		skip: dataQuery == null
	});

	const onlyData = useMemo(() => {

		if(!data?.data){
			return null;
		}

		const result = dataQuery?.getResultData(data.data);

		if (!result?.length) {
			return null;
		}

		return result[0][props.edge.name];
	}, [data.data]);

	return useMemo(() => ({
		data: onlyData, 
		loading: data.loading
	}), [onlyData]);
}