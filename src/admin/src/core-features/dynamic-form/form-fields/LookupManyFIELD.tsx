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
import { ContentManagerEntryDialogContent } from '@/features/content-manager/content-manager-new-item';
import { FormOpenMode } from '../form-field';
import { Options, Option } from '@/core-features/dynamic-form/form-field';
import { AutocompleteBaseFieldFormComponent, AutocompleteFieldFormComponentProps } from './AutocompleteBaseField';
import { ActionList } from '@/shared/components/action-list';
import { ActionListItem } from '@/shared/components/action-list-item';
import PlusCircleIcon from '@heroicons/react/24/outline/PlusCircleIcon';

interface RelationManyFIELDFormComponentProps {
    name: string;
    label: string;
    placeholder?: string;
    disabled?: boolean;
    required?: boolean;
    errors?: string;

    value: Array<Option<string, OptionTag>>;
    onChange: (event: any) => void;
    onBlur: (event: any) => void;

    entityName: string;
    editId?: string;
    edge: Edge;
}

export function LookupManyFIELDFormComponent(props: RelationManyFIELDFormComponentProps) {
    const { name, label, placeholder, disabled, errors, entityName, edge, value, onChange, onBlur, editId } = props;

    const myDialogContext = useMyDialogContext();
    const savedValue = useGetEdgeValue(entityName, editId ?? null, edge);

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

    const displayPropertyName = useMemo(() => fullEntity?.displayField?.name ?? 'name', [fullEntity?.displayField?.name]);


    const savedValueItems = useMemo((): Option<string, OptionTag>[] => {
        return savedValue?.data.map((i: any): Option<string, OptionTag> => ({
            ...i,
            tag: 'saved',
        })) ?? []
    }, [savedValue, displayPropertyName]);

    const connectedValueItems = useMemo((): Option<string, OptionTag>[] => {
        return value?.filter(i => i.tag === 'connect') ?? []
    }, [value]);

    const createdValueItems = useMemo((): Option<string, OptionTag>[] => {
        return value?.filter(i => i.tag === 'create') ?? []
    }, [value]);

    const disconnectedValueItems = useMemo((): Option<string, OptionTag>[] => {
        return value?.filter(i => i.tag === 'disconnect') ?? []
    }, [value]);

    const computedDisplayValue = useMemo((): Option<string, OptionTag>[] => {

        const disconnectIds = disconnectedValueItems.map(i => i.value);

        const result = [
            ...savedValueItems.filter((i: any) => disconnectIds.includes(i.value) === false),
            ...connectedValueItems,
            ...createdValueItems
        ]

        return result;
    }, [savedValueItems, connectedValueItems, createdValueItems, disconnectedValueItems]);


    const onChangedValue = useCallback((event: {
        reason: AutocompleteChangeReason;
        option: Option<string, OptionTag> | null;
    }) => {

        switch (event.reason) {
            case 'selectOption': {
                // check if selected item is in disconnect list
                if (disconnectedValueItems.find(i => i.value === event.option?.value)) {

                    const newValue = value.filter(i => (i.tag === 'disconnect' && i.value === event.option?.value) == false);
                    onChange({
                        target: {
                            name: name,
                            value: newValue
                        }
                    });
                    return;
                }

                // check if selected item is in saved list
                if (savedValueItems.find(i => i.value === event.option?.value)) {
                    return; // do nothing, it is already saved
                }

                const newValue = [
                    ...(value ?? []),
                    {
                        label: event.option?.label,
                        tag: 'connect',
                        value: event.option?.value
                    }]

                onChange({
                    target: {
                        name: name,
                        value: newValue
                    }
                });
            }
                break;

            case 'removeOption': {

                if (event.option?.tag === 'saved') {
                    const newItem = {
                        label: event.option.label,
                        tag: 'disconnect',
                        value: event.option.value
                    }

                    onChange({
                        target: {
                            name: name,
                            value: [...value ?? [], newItem]
                        }
                    });
                    return;
                }

                const newValue = value?.filter((i: any) => i.value !== event!.option!.value);
                onChange({
                    target: {
                        name: name,
                        value: newValue
                    }
                });

                break;
            }
            case 'clear': {

                const newValue = savedValueItems.map((i: any) => ({
                    label: i.label,
                    tag: 'disconnect',
                    value: i.value
                }))
                onChange({
                    target: {
                        name: name,
                        value: newValue
                    }
                });
                break;
            }
        }
    }, [onChange, savedValueItems, disconnectedValueItems]);

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

const useGetEdgeValue = (entityName: string, entryId: string | null, edge: Edge) => {

    const fullEntity = useFullEntity({ entityName: edge.relatedEntity.name });
    const contentManagerSearch = useContentManagerSearch({
        initialFilter: {
            id: {
                operator: eqStringFoldOperator.name,
                value: entryId
            }
        }
    });

    const contentManagerStore = useContentManagerStore({
        entityName: entityName,

        page: contentManagerSearch.state.page,
        rowsPerPage: contentManagerSearch.state.rowsPerPage,
        sortBy: contentManagerSearch.state.sortBy,
        order: contentManagerSearch.state.order,
        lazy: entryId == null,
        edges: useMemo(() => {
            if (!fullEntity?.displayField) {
                return undefined;
            }

            return [{
                name: edge.name,
                fields: ['id', fullEntity.displayField.name],
            }]
        }, [fullEntity?.displayField]),

        filters: useMemo(() => getAdvancedFiltersFromGridFilter(contentManagerSearch.state.filter), [contentManagerSearch.state.filter]),
    });

    const gridData = useMemo(() => {
        if (!contentManagerStore.state.data?.length) {
            return undefined;
        }

        return contentManagerStore.state.data[0][edge.name];
    }, [contentManagerStore.state.data]);

    return useMemo(() => ({
        data: gridData?.map((i: any): Option => {

            const label = i[fullEntity!.displayField!.name];

            return {
                label: label || `Id: ${i['id']}`,
                value: i.id,
            }
        }) ?? [],
        ids: gridData?.map((i: any) => i.id) ?? [],
    }), [gridData]);
}


export type OptionTag = 'saved' | 'connect' | 'create' | 'disconnect';