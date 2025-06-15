import { Edge } from "@/lib/apollo/graphql.entities";
import { useFullEntity } from "./use-entities";
import { useContentManagerSearch } from "./use-content-manager-search";
import { eqStringFoldOperator } from "@/core-features/dynamic-filter/filter-operators";
import { useContentManagerStore } from "./use-content-manager-store";
import { useMemo } from "react";
import { getAdvancedFiltersFromGridFilter } from "@/lib/utils/get-advanced-filters-from-grid-filters";

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