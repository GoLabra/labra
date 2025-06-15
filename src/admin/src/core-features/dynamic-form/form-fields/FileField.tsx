import { useMyDialogContext } from "@/core-features/dynamic-dialog/src/use-my-dialog-context";
import { useFormDynamicContext } from "../dynamic-form";
import { useFormContext } from "react-hook-form";
import { useLiteController } from "../lite-controller";
import { Edge } from "@/lib/apollo/graphql.entities";
import { UploadFilesBaseField } from "./UploadFilesBaseField";
import { DropEvent, FileRejection, FileWithPath } from "react-dropzone";
import { useCallback, useMemo } from "react";
import { EdgeStatus } from "@/lib/utils/edge-status";
import { useRelationDiff } from "@/features/content-manager/use-relation-diff";
import { createId } from "@paralleldrive/cuid2";
import { useGetEdgeValue } from "@/hooks/use-get-edge-value";
import { FileWithContent } from "@/shared/components/file-thumbnail";

export type FileDiffWrapper = {
	id?: string;
	file: FileWithPath | FileWithContent | string;
	status: EdgeStatus;
}

interface RelationManyFIELDFormComponentProps {
	name: string;
	label: string;
	placeholder?: string;
	disabled?: boolean;
	required?: boolean;
	errors?: string;

	value: FileDiffWrapper[];
	onChange: (event: any) => void;
	onBlur: (event: any) => void;

	entityName: string;
	editId?: string;
	edge: Edge;

	maxFiles?: number;
}

export function FileFieldFormComponent(props: RelationManyFIELDFormComponentProps) {
	const { name, label, placeholder, disabled, errors, entityName, edge, value, onChange, onBlur, editId } = props;

	const edgeValue = useGetEdgeValue<any[] | any>({
		entityName, 
		entryId: props.editId, 
		edge,
		fields: 'allfields'
	});

	const saved = useMemo((): FileDiffWrapper[] => {

		if(!edgeValue.data){
			return [];
		}

		const data = Array.isArray(edgeValue.data) ? edgeValue.data : [edgeValue.data];

		return data?.map(i => ({
			id: i.id,
			file: {
				name: i.name,
				content: i.content, //'https://content.app-sources.com/s/34322885479675261/uploads/Images/Wales-0386838.JPG',//i.name,
			},
			status: 'saved'
		})) ?? []

	}, [edgeValue.data]);


	const relationDiff = useRelationDiff<FileDiffWrapper>({ saved, changedArray: value });

	const files = useMemo(() => relationDiff.final.map(i => i.file), [relationDiff.final]);

	const onRemove = useCallback((file: File | string) => {
		onChange({
			target: {
				name: name,
				value: relationDiff.remove(i => {
					debugger;
					return i.file === file})
			}
		});
	}, [onChange, relationDiff.remove]);

	const onDrop = useCallback((acceptedFiles: File[], fileRejections: FileRejection[], event: DropEvent) => {
		if (!acceptedFiles.length) {
			return [];
		}

		debugger;
		onChange({
			target: {
				name: name,
				value: [
					...relationDiff.final,
					...acceptedFiles.map(i => ({
						id: createId(),
						file: i,
						status: 'create'
					}))
				]
			}
		});
	}, [value, relationDiff.final, onChange]);

	return (
		<>
			<UploadFilesBaseField
				name={name}
				label={label}
				placeholder={placeholder}
				value={files}
				disabled={disabled}
				maxFiles={props.maxFiles}
				errors={errors}
				onBlur={onBlur}
				onDrop={onDrop}
				onRemove={onRemove} />
		</>
	);
}

interface FileFormFieldProps {
	name: string;
	placeholder?: string;
	label: string;
	disabled?: boolean;
	hide?: boolean;
	required?: boolean;
	
	maxFiles?: number;
	entityName: string;
	edge: Edge;
}
export function FileFormField(props: FileFormFieldProps) {

	useFormDynamicContext(props.name, { disabled: props.disabled });
	const myDialogContext = useMyDialogContext();
	const formContext = useFormContext();
	const formControllerHandler = useLiteController({ name: props.name, control: formContext.control, disabled: props.disabled });

	if (props.hide) {
		return null;
	}

	return (<FileFieldFormComponent
		label={props.label}
		placeholder={props.placeholder}
		required={props.required}
		errors={formContext.formState.errors[props.name]?.message as string}
		editId={myDialogContext.editId}
		entityName={props.entityName}
		edge={props.edge}
		maxFiles={props.maxFiles}
		{...formControllerHandler}
	/>)
}