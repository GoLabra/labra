import { EdgeStatus } from "@/lib/utils/edge-status";
import { useCallback, useMemo } from "react";

export type RelationDiffItem<T = string> = {
	id: T
	status: EdgeStatus
}

interface RelationDiffProps<T extends RelationDiffItem> {
	saved: T[]
	changedArray: T[]
}
export const useRelationDiff = <T extends RelationDiffItem = any>(props: RelationDiffProps<T>) => {

	const showingItems = useMemo((): T[] => {
		const all = [...props.saved, ...props.changedArray ?? []];
		return all.filter(i => {
			if (i.status == 'delete') {
				return false;
			}

			if (i.status == 'disconnect') {
				return false;
			}

			if (i.status == 'saved') {
				if (all.find(j => j.id == i.id && (j.status == 'delete' || j.status == 'disconnect'))) {
					return false;
				}
			}

			return true;
		});
	}, [props.saved, props.changedArray]);

	const connect = useCallback((item: Omit<T, 'status'>) => {

		// check if selected item is in disconnect list
		const itemToRemove = props.changedArray?.find(i => i.id == item.id);
		if(itemToRemove){
			if(['disconnect', 'delete', 'unset'].includes(itemToRemove?.status)){
				return props.changedArray.filter(i => i != itemToRemove);		
			}
			return;
		}

		return [...props.changedArray ?? [],
					{
						...item,
						status: 'connect'
					}
				];
	
	}, [props.saved, props.changedArray]);

	// const remove = useCallback((find: (value: T, index: number, obj: T[]) => unknown) => {
	const remove = useCallback((id: string) => {

		const itemToRemove = showingItems.find(i => i.id == id);

		if(!itemToRemove){
			return;
		}
		// const savedItemToBeRemoved = props.saved.find(i => i.id == id);
		
		if (itemToRemove.status == 'saved') {
			return [
				...props.changedArray ?? [],
 				{
					...itemToRemove,
					status: 'delete'
				}
			]
		}

		return props.changedArray.filter(i => i != itemToRemove);


		// const itemToBeRemoved = props.changedArray.find(i => i.id == id);
		// if (!itemToBeRemoved) {
		// 	return props.changedArray;
		// }

		// return props.changedArray.filter(i => i != itemToBeRemoved);
	}, [props.saved, props.changedArray]);

	const disconnectAll = useCallback(() => {
		return props.saved.map(i => ({
			...i,
			status: 'disconnect'
		}));

	}, [props.saved]);
	
	return useMemo(() => ({
		showingItems,
		connect,
		remove,
		disconnectAll
	}), [showingItems, remove]);
}
