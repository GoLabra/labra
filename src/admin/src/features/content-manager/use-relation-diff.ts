import { EdgeStatus } from "@/lib/utils/edge-status";
import { useCallback, useMemo } from "react";

export type RelationDiffItem<T = string> = {
	id?: T
	status: EdgeStatus
}

interface RelationDiffProps<T extends RelationDiffItem> {
	saved: T[]
	changedArray: T[]
}
export const useRelationDiff = <T extends RelationDiffItem = any>(props: RelationDiffProps<T>) => {

	const final = useMemo((): T[] => {
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

	const remove = useCallback((find: (value: T, index: number, obj: T[]) => unknown) => {

		const savedItemToBeRemoved = props.saved.find(find);
		
		if (savedItemToBeRemoved) {
			return [
				...props.changedArray ?? [],
 				{
					...savedItemToBeRemoved,
					status: 'delete'
				}
			]
		}

		const itemToBeRemoved = props.changedArray.find(find);
		if (!itemToBeRemoved) {
			return props.changedArray;
		}

		//if (itemToBeRemoved.status == 'create' || itemToBeRemoved.status == 'connect') {
			return props.changedArray.filter(i => i != itemToBeRemoved);
		// }

		// return props.changedArray.map(i => {
		// 	if(i !== itemToBeRemoved){
		// 		return i;
		// 	}

		// 	return {
		// 		...i,
		// 		status: 'delete'
		// 	}
		// }); 
		
	}, [props.saved, props.changedArray]);

	return useMemo(() => ({
		final,
		remove
	}), [final, remove]);
}
