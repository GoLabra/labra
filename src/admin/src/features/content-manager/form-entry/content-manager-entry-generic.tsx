import { ChainDialogContentRef } from "@/core-features/dynamic-dialog/src/dynamic-dialog-types";
import { forwardRef, useCallback, useImperativeHandle, useMemo } from "react";
import { ContentManagerEntryDialogContentProps } from "../content-manager-entry-form";
import { DynamicDialogFooter } from "@/core-features/dynamic-dialog/src/use-dynamic-dialog-footer";
import { DynamicDialogHeader } from "@/core-features/dynamic-dialog/src/use-dynamic-dialog-header";
import { useMyDialogContext } from "@/core-features/dynamic-dialog/src/use-my-dialog-context";
import { FormOpenMode } from "@/core-features/dynamic-form/form-field";
import { useFullEntity } from "@/hooks/use-entities";
import { BasicAuditTrail } from "@/shared/components/basic-audit-trail";
import { Id } from "@/shared/components/id";
import { zodResolver } from "@hookform/resolvers/zod";
import { DialogContent, Stack, Button } from "@mui/material";
import { useForm } from "react-hook-form";
import { useContentManagerFormSchema } from "../use-content-manager-form-schema";
import { Form } from "@/core-features/dynamic-form2/dynamic-form";

export const ContentManagerEntryGeneric = forwardRef<ChainDialogContentRef, ContentManagerEntryDialogContentProps>((props, ref) => {

	const myDialogContext = useMyDialogContext();
	const fullEntity = useFullEntity({ entityName: props.entityName });
	const formSchema = useContentManagerFormSchema(fullEntity);

	const id = useMemo(() => {
		return props.defaultValue?.id;
	}, [props.defaultValue]);

	const displayValue = useMemo(() => {
		if (!fullEntity) {
			return null;
		}

		const displayFieldProperty = fullEntity?.displayField?.name;
		if (!displayFieldProperty) {
			return null;
		}

		if (displayFieldProperty == 'id') {
			return null;
		}

		return props.defaultValue?.[displayFieldProperty];
	}, [fullEntity, props.defaultValue]);

	const upperValueAsFieldValues = useMemo(() => {
		return Object.fromEntries(
			Object.entries(myDialogContext.upperResults ?? {}).map(([key, value]) => {
				return [key, value];
			}))
	}, [myDialogContext.upperResults]);

	const defaultValue = useMemo(() => {
		return {
			...formSchema.convertFromRawValue(formSchema.schemaDefaultValue),
			...formSchema.convertFromRawValue(props.defaultValue),
			//...formSchema.convertFromRawValue(upperValueAsFieldValues
		};

	}, [props.defaultValue, upperValueAsFieldValues, formSchema.convertFromRawValue]);

	const formMethods = useForm({
		resolver: zodResolver(formSchema.schema),
		mode: 'all',
		defaultValues: defaultValue
	});

	const onSave = useCallback((formData: any) => {
		myDialogContext.closeWithResults(formData);
	}, [myDialogContext.closeWithResults]);

	const header = useMemo(() => {
		const fEntityName = fullEntity?.caption ?? props.entityName;
		if (!displayValue) {
			return fEntityName;
		}
		return `${fEntityName}: ${displayValue}`;
	}, [displayValue, fullEntity?.caption, props.entityName]);

	useImperativeHandle(ref, () => ({
		enterPressed: () => {
			formMethods.handleSubmit(onSave)();
		}
	}));

	return (
		<>
			<DynamicDialogHeader>
				{header}
			</DynamicDialogHeader>
			<DialogContent>
				<Stack
					gap={2}>

					{!!id && <Id value={id} rootProps={{
						marginLeft: 'auto',
					}} />}

					<Form methods={formMethods} onSubmit={formMethods.handleSubmit(console.log)} >
						<Stack gap={1.5} >
							{formSchema.fields}
						</Stack>
					</Form>

					{myDialogContext.openMode !== FormOpenMode.New && <BasicAuditTrail defaultValues={props.defaultValue} />}

				</Stack>
			</DialogContent>
			<DynamicDialogFooter>
				<Stack
					direction="row"
					gap={1}>
					<Button
						color="primary"
						variant="contained"
						onClick={formMethods.handleSubmit(onSave)}
					>Save</Button>
				</Stack>
			</DynamicDialogFooter>
		</>
	)
});

ContentManagerEntryGeneric.displayName = 'ContentManagerEntryGeneric';