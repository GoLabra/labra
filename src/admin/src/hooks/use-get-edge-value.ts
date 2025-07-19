import { Edge } from "@/lib/apollo/graphql.entities";
import { useFullEntity } from "./use-entities";
import { useContentManagerSearch } from "./use-content-manager-search";
import { eqStringFoldOperator } from "@/core-features/dynamic-filter/filter-operators";
import { useContentManagerStore } from "./use-content-manager-store";
import { useMemo } from "react";
import { getAdvancedFiltersFromGridFilter } from "@/lib/utils/get-advanced-filters-from-grid-filters";
import { gql, useQuery } from "@apollo/client";
import { fileIsImage, fileTypeByUrl } from "@/shared/components/file-thumbnail";
import { EdgeRequest } from "@/lib/apollo/builders/gqlQueryBuilder";
import { ADMIN_CONTEXT } from "@/lib/apollo/apolloWrapper";

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

	const edgeValueData = useGetEdgeValue<any[] | any>({
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


interface UseGetEdgeValueParams {
	entityName: string;
	entryId: string | null | undefined;
	edge: Edge;
	fields: 'allfields' | 'iddisplay' | string[];
	edges?: EdgeRequest[];
}
export const useGetEdgeValue = <T = any>(props: UseGetEdgeValueParams) => {

	const rootEntity = useFullEntity({ entityName: props.entityName });
	const edgeEntity = useFullEntity({ entityName: props.edge.relatedEntity.name });
	const contentManagerSearch = useContentManagerSearch({
		initialFilter: {
			id: {
				operator: eqStringFoldOperator.name,
				value: props.entryId
			}
		}
	});

	const fields = useMemo(() => {
		if (props.fields === 'allfields') {
			return edgeEntity?.fields?.map(i => i.name) ?? undefined;
		}

		if (props.fields === 'iddisplay') {

			if (!edgeEntity?.displayField) {
				return undefined;
			}
			
			return ['id', edgeEntity?.displayField?.name];
		}

		return props.fields;
	}, [edgeEntity?.displayField, edgeEntity?.fields, props.fields]);

	const edges = useMemo(() => {
		return props.edges ?? [];
	}, [props.edges]);

	const contentManagerStore = useContentManagerStore({
		entityOwner: rootEntity?.owner,
		entityName: props.entityName,

		page: contentManagerSearch.state.page,
		rowsPerPage: contentManagerSearch.state.rowsPerPage,
		sortBy: contentManagerSearch.state.sortBy,
		order: contentManagerSearch.state.order,
		skip: props.entryId == null || (!props.fields.length && !edges.length),
		edges: useMemo(() => {
			if (!fields) {
				return undefined;
			}

			if (!fields.length) {
				return undefined;
			}

			return [{
				name: props.edge.name,
				fields: fields,
				edges: edges
			}]
		}, [edgeEntity?.displayField]),

		filters: useMemo(() => getAdvancedFiltersFromGridFilter(contentManagerSearch.state.filter), [contentManagerSearch.state.filter]),
	});

	const data: T = useMemo(() => {
		if (!contentManagerStore.state.data?.length) {
			return null;
		}

		return contentManagerStore.state.data[0][props.edge.name];
	}, [contentManagerStore.state.data]);

	return useMemo(() => ({
		data: data, 
	}), [contentManagerStore.state.data]);
}