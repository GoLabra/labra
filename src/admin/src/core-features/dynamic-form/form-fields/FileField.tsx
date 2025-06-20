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
import { useEntityFiles, useGetEdgeValue } from "@/hooks/use-get-edge-value";
import { FileData } from "@/shared/components/file-thumbnail";

export type FileDiffWrapper = {
	id?: string;
	file: FileWithPath | FileData | string;
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

	const edgeFiles = useEntityFiles({
		entityName, 
		entryId: props.editId, 
		edge
	});

	const saved = useMemo((): FileDiffWrapper[] => {
		return edgeFiles?.map(i => ({
			id: i.id,
			file: {
				name: i.name,
				mimeType: '',
				preview: i.content,
			},
			status: 'saved'
		})) ?? []

	}, [edgeFiles]);


	const relationDiff = useRelationDiff<FileDiffWrapper>({ saved, changedArray: value });

	const files = useMemo(() => relationDiff.showingItems.map(i => i.file), [relationDiff.showingItems]);

	const onRemove = useCallback((file: File | string) => {
		onChange({
			target: {
				name: name,
				value: relationDiff.remove(i => i.file === file)
			}
		});
	}, [onChange, relationDiff.remove]);

	const onDrop = useCallback((acceptedFiles: File[], fileRejections: FileRejection[], event: DropEvent) => {
		if (!acceptedFiles.length) {
			return [];
		}

		onChange({
			target: {
				name: name,
				value: [
					...value ?? [],
					...acceptedFiles.map(i => ({
						id: createId(),
						file: i,
						status: 'create'
					}))
				]
			}
		});
	}, [value, onChange]);

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