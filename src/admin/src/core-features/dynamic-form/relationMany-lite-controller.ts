import { GenericEvent } from "@/lib/utils/event";
import { useCallback, useMemo, useRef } from "react";
import { Control, useWatch } from "react-hook-form";
import { useLiteController } from "./lite-controller";

export const RelationInfoType = '_RELATION_INFO_';
export type RelationInfo = {
	savedCount: number;
	type?: typeof RelationInfoType;
} 

interface useLiteControllerProps {
    name: string;
    control: Control;
    disabled?: boolean;
}
export const useRelationManyLiteController = <T = any>(props: useLiteControllerProps) => {

    const liteController = useLiteController<any[]>({ name: props.name, control: props.control, disabled: props.disabled });
	const relationInfoRef = useRef<RelationInfo>();

	const addRelationInfo = useCallback((relationInfo: RelationInfo) => {
		relationInfoRef.current = {
			...relationInfo,
			type: RelationInfoType
		};

		liteController.onChange({
			target: {
				name: props.name,
				value: [
					...(liteController.value ?? []).filter((i: RelationInfo) => i.type !== RelationInfoType),
					relationInfoRef.current
				]
			}
		});
	}, [liteController.onChange]);

	const value = useMemo(() => {
		return liteController.value?.filter((i: RelationInfo) => i.type !== RelationInfoType);
	}, [liteController.value]);

	const onChange = useCallback((event: GenericEvent<T> | any) => {
		liteController.onChange({
			target: {
				name: props.name,
				value: [
					...event.target.value ?? [],
					...(relationInfoRef.current ? [relationInfoRef.current] : [])
				]
			}
		});
	}, [liteController.onChange]);

	return useMemo(() => ({
		onBlur: liteController.onBlur,
		disabled: liteController.disabled,
		name: liteController.name,
		value,
		onChange,
		addRelationInfo
	}), [liteController.onBlur, liteController.disabled, liteController.name, value, onChange, addRelationInfo]);
}
