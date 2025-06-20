import { Edge } from "@/lib/apollo/graphql.entities";
import { useFullEntity } from "./use-entities";
import { useContentManagerSearch } from "./use-content-manager-search";
import { eqStringFoldOperator } from "@/core-features/dynamic-filter/filter-operators";
import { useContentManagerStore } from "./use-content-manager-store";
import { useMemo } from "react";
import { getAdvancedFiltersFromGridFilter } from "@/lib/utils/get-advanced-filters-from-grid-filters";
import { gql, useQuery } from "@apollo/client";
import { fileIsImage, fileTypeByUrl } from "@/shared/components/file-thumbnail";

const GET_FILES_CONTENT = gql`query files($or: [FileWhereInput!]) {
	files(where: { or: $or }) {
		id
		content
	}
}`


interface UseEntityFilesParams {
	entityName: string;
	entryId: string | null | undefined;
	edge: Edge;
}
export const useEntityFiles = (props: UseEntityFilesParams) => {

	const edgeValueData = useGetEdgeValue<any[] | any>({
		entityName: props.entityName, 
		entryId: props.entryId, 
		edge: props.edge,
		fields: ['id', 'name']
	});

	const files = useMemo(() => {
		if (!edgeValueData.data) {
			return [];
		}

		if(!Array.isArray(edgeValueData.data)){
			return [edgeValueData.data]
		}

		return edgeValueData.data ?? [];
	}, [edgeValueData.data]);

	const previewIds = useMemo(() => {
		return files.filter(i => fileIsImage(fileTypeByUrl(i.name)))
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
		skip: previewIds.length == 0

	});

	// console.log('filesContent', filesContent);

	return useMemo((): any[] => {
		
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

}


interface UseGetEdgeValueParams {
	entityName: string;
	entryId: string | null | undefined;
	edge: Edge;
	fields: 'allfields' | 'iddisplay' | string[] 
}
export const useGetEdgeValue = <T = any>(props: UseGetEdgeValueParams) => {

	const fullEntity = useFullEntity({ entityName: props.edge.relatedEntity.name });
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
			return fullEntity?.fields?.map(i => i.name) ?? undefined;
		}

		if (props.fields === 'iddisplay') {

			if (!fullEntity?.displayField) {
				return undefined;
			}
			
			return ['id', fullEntity?.displayField?.name];
		}

		return props.fields;
	}, [fullEntity?.displayField, fullEntity?.fields, props.fields]);

	const contentManagerStore = useContentManagerStore({
		entityName: props.entityName,

		page: contentManagerSearch.state.page,
		rowsPerPage: contentManagerSearch.state.rowsPerPage,
		sortBy: contentManagerSearch.state.sortBy,
		order: contentManagerSearch.state.order,
		lazy: props.entryId == null,
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
			}]
		}, [fullEntity?.displayField]),

		filters: useMemo(() => getAdvancedFiltersFromGridFilter(contentManagerSearch.state.filter), [contentManagerSearch.state.filter]),
	});

	const data: T = useMemo(() => {
		console.log(contentManagerStore.state.data);
		if (!contentManagerStore.state.data?.length) {
			return null;
		}

		return contentManagerStore.state.data[0][props.edge.name];
	}, [contentManagerStore.state.data]);

	return useMemo(() => ({
		data: data, 
	}), [data]);
}