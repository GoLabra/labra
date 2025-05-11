"use client";

import { useFormContext } from 'react-hook-form';
import { AutocompleteBaseFieldFormComponent, AutocompleteFieldFormComponentProps } from './AutocompleteBaseField';
import { FormOpenMode } from '@/core-features/dynamic-form/form-field';
import { AutocompleteChangeReason, SvgIcon } from '@mui/material';
import PlusCircleIcon from '@heroicons/react/24/outline/PlusCircleIcon';
import { useCallback, useContext, useMemo } from 'react';
import { ContentManagerEntryDialogContent } from '@/features/content-manager/content-manager-new-item';
import { ActionList } from '@/shared/components/action-list';
import { ActionListItem } from '@/shared/components/action-list-item';
import { useFullEntity } from '@/hooks/use-entities';
import { useContentManagerSearch } from '@/hooks/use-content-manager-search';
import { useContentManagerStore } from '@/hooks/use-content-manager-store';
import { useMyDialogContext } from '@/core-features/dynamic-dialog/src/use-my-dialog-context';
import { useFormDynamicContext } from '@/core-features/dynamic-form2/dynamic-form';
import { useLiteController } from '../lite-controller';
import { Edge } from '@/lib/apollo/graphql.entities';
import { getAdvancedFiltersFromGridFilter } from '@/lib/utils/get-advanced-filters-from-grid-filters';
import { eqStringFoldOperator } from '@/core-features/dynamic-filter/filter-operators';
import { getAdvancedFiltersFromQuery } from '@/lib/utils/get-filters-from-query';
import { Options, Option } from '@/core-features/dynamic-form/form-field';

interface RelationOneFIELDFormComponentProps {
    name: string;
    label: string;
    placeholder?: string;
    disabled?: boolean;
    required?: boolean;
    errors?: string;

    value: Option<string, OptionTag>;
    onChange: (event: any) => void;
    onBlur: (event: any) => void;

    entityName: string;
    editId?: string;
    edge: Edge;
}

export function LookupOneFIELDFormComponent(props: RelationOneFIELDFormComponentProps) {
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

    const computedValue = useMemo((): Option<string, OptionTag> | undefined => {
        if (value) {
            if(value.tag === 'unset'){
                return undefined;
            }
            return value;
        }

        if(!savedValue?.data){
            return undefined;
        }

        if (savedValue) {
            return {
                ...savedValue.data!,
                tag: 'saved'
            };
        }

        return undefined;
    }, [value, savedValue, fullEntity?.displayField]);

    const onChangedValue = useCallback((event: {
        reason: AutocompleteChangeReason;
        option: Option<string, OptionTag> | null;
    }) => {

        switch (event.reason) {
            case 'selectOption': {
                onChange({
                    target: {
                        name: name,
                        value: {
                            ...event.option,
                            tag: 'connect'
                        }
                    }
                });
                break;
            }

            case 'removeOption':
            case 'clear': {

                if (!savedValue?.data) {
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
                            label: savedValue?.data.label,
                            value: savedValue?.data.value,
                            tag: 'unset'
                        }
                    }
                });
                break;
            }
        }

    }, [onChange, savedValue]);

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
                        value:  {
                            label: upperResults[displayPropertyName],
                            tag: 'create',
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
    const formControllerHandler = useLiteController({ name: props.name, control: formContext.control, disabled: props.disabled });
    
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

    return useMemo(() => {

        if(!gridData){
            return {
                data: undefined,
                id: undefined,
            }
        }

        return {
            data: {
                label: gridData[fullEntity!.displayField!.name] || `Id: ${gridData['id']}`,
                value: gridData.id,
            },
            id: gridData.id,
        }
    }, [gridData]);
}

export type OptionTag = 'saved' | 'connect' | 'create' | 'unset';