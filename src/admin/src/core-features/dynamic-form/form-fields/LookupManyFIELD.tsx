"use client";

import { Control, FieldError, FieldErrorsImpl, Merge, useFormContext, UseFormRegister } from 'react-hook-form';
import { gql } from "@apollo/client";
import { AutocompleteChangeReason, Box, Chip, IconButton, List, ListItemButton, ListItemText, Stack, SvgIcon, TextField } from '@mui/material';
import { useFormDynamicContext } from '@/core-features/dynamic-form2/dynamic-form';
import { useLiteController } from '../lite-controller';
import { useMyDialogContext } from '@/core-features/dynamic-dialog/src/use-my-dialog-context';
import { Edge } from '@/lib/apollo/graphql.entities';
import { useFullEntity } from '@/hooks/use-entities';
import { useContentManagerSearch } from '@/hooks/use-content-manager-search';
import { eqStringFoldOperator } from '@/core-features/dynamic-filter/filter-operators';
import { useContentManagerStore } from '@/hooks/use-content-manager-store';
import { getAdvancedFiltersFromGridFilter } from '@/lib/utils/get-advanced-filters-from-grid-filters';
import { useCallback, useMemo } from 'react';
import { getAdvancedFiltersFromQuery } from '@/lib/utils/get-filters-from-query';
import { ContentManagerEntryDialogContent } from '@/features/content-manager/content-manager-entry-form';
import { FormOpenMode } from '../form-field';
import { Options, Option } from '@/core-features/dynamic-form/form-field';
import { AutocompleteBaseFieldFormComponent, AutocompleteFieldFormComponentProps } from './AutocompleteBaseField';
import { ActionList } from '@/shared/components/action-list';
import { ActionListItem } from '@/shared/components/action-list-item';
import PlusCircleIcon from '@heroicons/react/24/outline/PlusCircleIcon';
import { EdgeStatus } from '@/lib/utils/edge-status';
import { useRelationDiff } from '@/features/content-manager/use-relation-diff';
import { OptionDiffWrapper } from './LookupOneFIELD';
import { useGetEdgeValue } from '@/hooks/use-get-edge-value';

interface RelationManyFIELDFormComponentProps {
    name: string;
    label: string;
    placeholder?: string;
    disabled?: boolean;
    required?: boolean;
    errors?: string;

    value: OptionDiffWrapper[];
    onChange: (event: any) => void;
    onBlur: (event: any) => void;

    entityName: string;
    editId?: string;
    edge: Edge;
}

export function LookupManyFIELDFormComponent(props: RelationManyFIELDFormComponentProps) {
    const { name, label, placeholder, disabled, errors, entityName, edge, value, onChange, onBlur, editId } = props;

    
    
    // SEARCH LOOKUP
    const fullEntity = useFullEntity({ entityName: edge.relatedEntity.name });
    const contentManagerSearch = useContentManagerSearch();
    const contentManagerStore = useContentManagerStore({
        entityName: edge.relatedEntity.name,
        page: contentManagerSearch.state.page,
        rowsPerPage: contentManagerSearch.state.rowsPerPage,
        sortBy: contentManagerSearch.state.sortBy,
        order: contentManagerSearch.state.order,
        fields: useMemo(() => {
            if (!fullEntity?.displayField) {
                return [];
            }

            return ['id', fullEntity.displayField.name];

        }, [fullEntity?.displayField]),

        orFilters: useMemo(() => getAdvancedFiltersFromQuery(contentManagerSearch.state.query, fullEntity?.fields ?? []), [contentManagerSearch.state.query, fullEntity?.fields])
    });

    const search = useCallback((searchValue: string) => {
        contentManagerSearch.handleQueryChange(searchValue);
    }, [contentManagerSearch]);
    // END SEARCH LOOKUP

	const myDialogContext = useMyDialogContext();
	const savedValue = useGetEdgeValue({
		entityName: entityName, 
		entryId: editId, 
		edge: props.edge,
		fields: 'iddisplay'
	});
    const displayPropertyName = useMemo(() => fullEntity?.displayField?.name ?? 'name', [fullEntity?.displayField?.name]);

    const savedValueItems = useMemo((): OptionDiffWrapper[] => {
        return savedValue?.data?.map((i: any) => ({
            id: i.id,
			label: i[displayPropertyName],
			value: i,
            status: 'saved',
        })) ?? []
    }, [savedValue.data]);

	const relationDiff = useRelationDiff<OptionDiffWrapper>({ saved: savedValueItems, changedArray: value });
	
	const computedDisplayValue = useMemo((): Option[] => relationDiff.showingItems.map(i => ({
		label: i.label,
		value: i.id
	})), [relationDiff.showingItems]);

    const onChangedValue = useCallback((event: {
        reason: AutocompleteChangeReason;
        option: Option<string> | null;
    }) => {

        switch (event.reason) {
            case 'selectOption': 
				onChange({
					target: {
						name: name,
						value: relationDiff.connect({
									id: event.option!.value,
									label: event.option!.label,
									value: event.option!.value
								})
					}
				});
                break;

            case 'removeOption': 
				onChange({
                    target: {
                        name: name,
                        value: relationDiff.remove(event.option!.value)
                    }
                });
                break;
            
            case 'clear': 
                onChange({
					target: {
						name: name,
						value: relationDiff.disconnectAll()
					}
				});
                break;
            
        }
    }, [onChange, savedValueItems, relationDiff.connect]);

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
                        value: [
                            ...(value ?? []),
                            {
                                label: upperResults[displayPropertyName],
                                tag: 'create',
                                value: upperResults
                            }
                        ]
                    }
                });

                return upperResults;
            })
    };

    return (
        <>
            <AutocompleteBaseFieldFormComponent
                multiple
                name={name}
                label={label}
                placeholder={placeholder}
                disabled={disabled}
                errors={errors}
                search={search}
                options={options}

                value={computedDisplayValue}
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
export function LookupManyFIELDFormField(props: FormFieldProps) {

    useFormDynamicContext(props.name, { disabled: props.disabled });
    const myDialogContext = useMyDialogContext();
    const formContext = useFormContext();
    const formControllerHandler = useLiteController({ name: props.name, control: formContext.control, disabled: props.disabled });

    if (props.hide) {
        return null;
    }

    return (<LookupManyFIELDFormComponent
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
