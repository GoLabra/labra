"use client";

import { useFormContext } from 'react-hook-form';
import { AutocompleteBaseFieldFormComponent, AutocompleteFieldFormComponentProps } from './AutocompleteBaseField';
import { FormOpenMode } from '@/core-features/dynamic-form/form-field';
import { AutocompleteChangeReason, SvgIcon } from '@mui/material';
import PlusCircleIcon from '@heroicons/react/24/outline/PlusCircleIcon';
import { useCallback, useContext, useMemo } from 'react';
import { ContentManagerEntryDialogContent } from '@/features/content-manager/content-manager-entry-form';
import { ActionList } from '@/shared/components/action-list';
import { ActionListItem } from '@/shared/components/action-list-item';
import { useFullEntity } from '@/hooks/use-entities';
import { useContentManagerSearch } from '@/hooks/use-content-manager-search';
import { useMyDialogContext } from '@/core-features/dynamic-dialog/src/use-my-dialog-context';
import { useFormDynamicContext } from '@/core-features/dynamic-form2/dynamic-form';
import { useLiteController } from '../lite-controller';
import { Edge } from '@/lib/apollo/graphql.entities';
import { Options, Option } from '@/core-features/dynamic-form/form-field';
import { EdgeStatus } from '@/lib/utils/edge-status';
import { nullable } from 'zod';
import { useLookupContentManagerStore } from '@/hooks/use-lookup-content-manager-store';
import { useRelationContentManagerStore } from '@/hooks/use-relation-content-manager-store';

export type OptionDiffWrapper = {
	id: string;
	label: string;
	value: string | any;
	status: EdgeStatus;
}

interface RelationOneFIELDFormComponentProps {
	name: string;
	label: string;
	placeholder?: string;
	disabled?: boolean;
	required?: boolean;
	errors?: string;

	value: OptionDiffWrapper;
	onChange: (event: any) => void;
	onBlur: (event: any) => void;

	entityName: string;
	editId?: string;
	edge: Edge;
}

export function LookupOneFIELDFormComponent(props: RelationOneFIELDFormComponentProps) {
	const { name, label, placeholder, disabled, errors, entityName, edge, value, onChange, onBlur, editId } = props;


	// SEARCH LOOKUP
	const fullEntity = useFullEntity({ entityName: edge.relatedEntity.name });
	const contentManagerSearch = useContentManagerSearch();
	const contentManagerStore = useLookupContentManagerStore({
		fullEntity: fullEntity,
		searchState: contentManagerSearch.state,
	});

	const search = useCallback((searchValue: string) => {
		contentManagerSearch.handleQueryChange(searchValue);
	}, [contentManagerSearch]);
	// END SEARCH LOOKUP


	const myDialogContext = useMyDialogContext();
	const displayPropertyName = useMemo(() => fullEntity?.displayField?.name ?? 'name', [fullEntity?.displayField?.name]);
	const edgeValue = useRelationContentManagerStore({
		entityName: entityName,
		entryId: editId,
		edge: edge,
		fields: 'iddisplay'
	});

	const saved = useMemo((): OptionDiffWrapper | null => {
		if (!fullEntity?.displayField) {
			return null;
		}

		if (!edgeValue.data) {
			return null;
		}

		return {
			id: edgeValue.data.id,
			value: edgeValue.data.id,
			label: edgeValue.data[displayPropertyName],
			status: 'saved'
		};
	}, [edgeValue, displayPropertyName]);

	const computedValue = useMemo((): OptionDiffWrapper | undefined => {
		if (value) {
			if (value.status === 'unset') {
				return undefined;
			}
			return value;
		}

		if (!saved) {
			return undefined;
		}

		if (saved) {
			return {
				...saved,
				status: 'saved'
			};
		}

		return undefined;
	}, [value, saved, fullEntity?.displayField]);

	const onChangedValue = useCallback((event: {
		reason: AutocompleteChangeReason;
		option: OptionDiffWrapper | null;
	}) => {

		switch (event.reason) {
			case 'selectOption': {
				onChange({
					target: {
						name: name,
						value: {
							...event.option,
							status: 'connect'
						}
					}
				});
				break;
			}

			case 'removeOption':
			case 'clear': {

				if (!saved) {
					onChange({
						target: {
							name: name,
							value: null
						}
					});
					return;
				}

				onChange({
					target: {
						name: name,
						value: {
							label: saved.label,
							value: saved.value,
							status: 'unset'
						}
					}
				});
				break;
			}
		}

	}, [onChange, saved]);

	const options = useMemo(() => {

		const result = contentManagerStore.state.data.map((i: any): Option<any> => ({
			label: i[displayPropertyName] || `Id: ${i['id']}`,
			value: i.id,
		}));

		return result;

	}, [contentManagerStore.state.data]);

	const addNew = () => {
		myDialogContext.addPopup(name, ContentManagerEntryDialogContent, { entityName: edge.relatedEntity.name }, FormOpenMode.New, undefined,
			(upperResults: any) => {

				onChange({
					target: {
						name: name,
						value: {
							label: upperResults[displayPropertyName],
							status: 'create',
							value: upperResults
						}
					}
				});

				return upperResults;
			})
	};

	return (
		<>
			<AutocompleteBaseFieldFormComponent
				name={name}
				label={label}
				placeholder={placeholder}
				disabled={disabled}
				errors={errors}
				search={search}
				options={options}

				value={computedValue}
				onChange={onChangedValue as AutocompleteFieldFormComponentProps['onChange']}
				onBlur={onBlur}

				customHeader={
					<ActionList>
						<ActionListItem
							onClick={() => addNew()}
							icon={(
								<SvgIcon fontSize="small">
									<PlusCircleIcon />
								</SvgIcon>
							)}
							aria-label="Add new entry"
							aria-haspopup="dialog"
							label="Add New"
						/>
					</ActionList>
				}
			/>
		</>
	);
}

interface FormFieldProps {
	name: string;
	placeholder?: string;
	label: string;
	disabled?: boolean;
	hide?: boolean;
	required?: boolean;

	entityName: string;
	edge: Edge;
}
export function LookupOneFIELDFormField(props: FormFieldProps) {
	useFormDynamicContext(props.name, { disabled: props.disabled });
	const myDialogContext = useMyDialogContext();
	const formContext = useFormContext();
	const formControllerHandler = useLiteController({ name: props.name, disabled: props.disabled });

	if (props.hide) {
		return null;
	}

	return (<LookupOneFIELDFormComponent
		label={props.label}
		placeholder={props.placeholder}
		required={props.required}
		errors={formContext.formState.errors[props.name]?.message as string}
		editId={myDialogContext.editId}
		entityName={props.entityName}
		edge={props.edge}
		{...formControllerHandler}
	/>)
}