import { Edge } from "@/lib/apollo/graphql.entities";
import { createRowDetailStore, createRowSelectionStore, RowDetailStore, useRowDetails } from "mosaic-data-table";
import { useCallback, useMemo } from "react";

export type EdgeNode = {
	entityName: string;
	edge: Edge;
	entryId: string;
}
export type RelationViewState = {
	edges: EdgeNode[];
}

export const useViewRelationStore = () => {
	const detailsStore = useMemo(() => createRowDetailStore<any>(), [])

	const close = useCallback((rootEntryId: string) => {
		detailsStore.clear(rootEntryId)
	}, [detailsStore]);

	const closeAll = useCallback(() => {
		detailsStore.clear();
	}, [detailsStore]);

	const addEdge = useCallback((rootEntryId: string, entityName: string, edge: Edge, entryId: string, addAsFirst: boolean = false) => {

		if (addAsFirst) {
			detailsStore.getExpansionInfo
			detailsStore.setParams(rootEntryId, {
				edges: [
					{
						entityName,
						edge,
						entryId
					}]
			}, true);

			return;
		}

		const edgeNodesInfo = detailsStore.getExpansionInfo(rootEntryId);
		const params = edgeNodesInfo.params ?? {};
		const edges = params?.edges ?? [];

		detailsStore.setParams(rootEntryId, {
			edges: [
				...edges,
				{
					entityName,
					edge,
					entryId
				}]
		}, true);

	}, [detailsStore]);

	const goBack = useCallback((rootEntryId: string) => {
		const edgeNodesInfo = detailsStore.getExpansionInfo(rootEntryId);
		const params = edgeNodesInfo.params ?? {};
		const edges = params?.edges ?? [];

		detailsStore.setParams(rootEntryId, {
				edges: [
					...edges.slice(0, -1)
				]
			},
			true
		);
	}, [detailsStore]);

	return useMemo(() => ({
		detailsStore,
		addEdge,
		goBack,
		close,
		closeAll
	}), [detailsStore, goBack, addEdge, close, closeAll]);
}


export const useOneViewRelationStore = (rootEntryId: string, viewRelationStore: ReturnType<typeof useViewRelationStore>) => {

	const detailStore = useRowDetails(viewRelationStore.detailsStore, rootEntryId);
	
	const edges = useMemo((): EdgeNode[] => {
		const params = detailStore.params ?? {};
		const edges = params?.edges ?? [];
		return edges;
	}, [rootEntryId, detailStore]);

	const visibleEdge = useMemo(() => {
		return edges.findLast(i => i);
	}, [edges]);

	const addEdge = useCallback((entityName: string, edge: Edge, entryId: string) => {
		viewRelationStore.addEdge(rootEntryId, entityName, edge, entryId); 
	}, [viewRelationStore, rootEntryId]);

	const goBack = useCallback(() => {
		viewRelationStore.goBack(rootEntryId);
	}, [viewRelationStore, rootEntryId]);

	const close = useCallback(() => {
		viewRelationStore.close(rootEntryId);
	}, [viewRelationStore, rootEntryId]);

	return useMemo(() => ({
		edges,
		visibleEdge,
		addEdge,
		goBack,
		close
	}), [edges, visibleEdge, addEdge, goBack, close]);
}
