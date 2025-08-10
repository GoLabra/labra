import { useMyDialogContext } from "@/core-features/dynamic-dialog/src/use-my-dialog-context";
import { useFormDynamicContext } from "../dynamic-form";
import { useFormContext } from "react-hook-form";
import { useLiteController } from "../lite-controller";
import { Edge } from "@/lib/apollo/graphql.entities";
import { UploadFilesBaseField } from "./UploadFilesBaseField";
import { DropEvent, FileRejection, FileWithPath } from "react-dropzone";
import { useCallback, useEffect, useMemo } from "react";
import { EdgeStatus } from "@/lib/utils/edge-status";
import { useRelationDiff } from "@/features/content-manager/use-relation-diff";
import { createId } from "@paralleldrive/cuid2";
import { LFile, useEntityFiles } from "@/hooks/use-relation-content-manager-store";
import { FileData, fileIsImage, fileTypeByUrl } from "@/shared/components/file-thumbnail";
import { RelationInfo, useRelationManyLiteController } from "../relationMany-lite-controller";
import { createGenericEvent } from "@/lib/utils/event";

export type FileDiffWrapper = {
	id: string;
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
	addRelationInfo?: (relationInfo: RelationInfo) => void;
}

export function FileFieldFormComponent(props: RelationManyFIELDFormComponentProps) {
	const { name, label, placeholder, disabled, errors, entityName, edge, value, onChange, onBlur, editId } = props;

	const edgeFiles = useEntityFiles({
		entityName, 
		entryId: props.editId, 
		edge,
		loadContent: useCallback((file: LFile) => {
			return fileIsImage(fileTypeByUrl(file.name));
		}, []),
	});

	useEffect(() => {
		props.addRelationInfo?.({
			savedCount: edgeFiles.length
		});
	}, [edgeFiles.length, props.addRelationInfo]);

	const savedValueItems = useMemo((): FileDiffWrapper[] => {
		return edgeFiles?.map(i => ({
			id: i.id,
			file: {
				caption: i.caption,
				name: i.name,
				size: i.size,
				mimeType: i.mimeType,
				...(i.content && {
					preview: `data:${i.mimeType};base64,${i.content}`,
				})
			},
			status: 'saved'
		})) ?? []
	}, [edgeFiles]);


	const relationDiff = useRelationDiff<FileDiffWrapper>({ saved: savedValueItems, changedArray: value });

	const files = useMemo(() => relationDiff.showingItems.map(i => i.file), [relationDiff.showingItems]);

	const onRemove = useCallback((file: File | FileData | string) => {

		const removeId = relationDiff.showingItems.find(i => i.file === file)?.id;
		if(!removeId){
			return;
		}

		onChange(createGenericEvent(props.name,
				relationDiff.remove(removeId))
		);
	}, [onChange, relationDiff.remove]);

	const onDrop = useCallback((acceptedFiles: File[], fileRejections: FileRejection[], event: DropEvent) => {
		if (!acceptedFiles.length) {
			return [];
		}

		onChange(createGenericEvent(props.name,
				[
					...value ?? [],
					...acceptedFiles.map(i => ({
						id: createId(),
						file: i,
						status: 'create'
					}))
				]
		));
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
	const formControllerHandler = useRelationManyLiteController({ name: props.name, disabled: props.disabled });

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