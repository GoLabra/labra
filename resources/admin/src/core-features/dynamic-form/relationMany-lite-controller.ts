import { createGenericEvent, GenericEvent } from "@/lib/utils/event";
import { useCallback, useMemo, useRef } from "react";
import { Control, useFormContext, useWatch } from "react-hook-form";
import { useLiteController } from "./lite-controller";

export const RelationInfoType = '_RELATION_INFO_';
export type RelationInfo = {
	savedCount: number;
	type?: typeof RelationInfoType;
} 

interface useLiteControllerProps {
    name: string;
    disabled?: boolean;
}
export const useRelationManyLiteController = <T = any>(props: useLiteControllerProps) => {

	const formContext = useFormContext();
    const liteController = useLiteController<any[]>({ name: props.name, disabled: props.disabled });
	const relationInfoRef = useRef<RelationInfo>();

	const addRelationInfo = useCallback((relationInfo: RelationInfo) => {
		relationInfoRef.current = {
			...relationInfo,
			type: RelationInfoType
		};

		formContext.setValue(props.name, [
					...(liteController.value ?? []).filter((i: RelationInfo) => i.type !== RelationInfoType),
					relationInfoRef.current
				],{
					shouldValidate: false,
					shouldDirty: false,
					shouldTouch: false,
				}
		);
	}, [liteController.onChange]);

	const value = useMemo(() => {
		return liteController.value?.filter((i: RelationInfo) => i.type !== RelationInfoType);
	}, [liteController.value]);

	const onChange = useCallback((event: GenericEvent<T> | any) => {
		liteController.onChange(createGenericEvent(props.name,
				[
					...event.target.value ?? [],
					...(relationInfoRef.current ? [relationInfoRef.current] : [])
				])
		);
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
