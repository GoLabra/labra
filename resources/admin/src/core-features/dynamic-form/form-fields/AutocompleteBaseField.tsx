"use client";

import { Autocomplete, AutocompleteChangeReason, AutocompletePopperSlotPropsOverrides, Button, Paper, Popper, PopperProps, SvgIcon, TextField, debounce, styled } from '@mui/material';
import { Control, FieldError, FieldErrorsImpl, Merge, UseFormRegister } from 'react-hook-form';
import Box from '@mui/material/Box';
import { useCallback } from 'react';
import Grid from '@mui/material/Grid';
import Typography from '@mui/material/Typography';
import { Options, Option } from '@/core-features/dynamic-form/form-field';
import React from 'react';
import { useLiteController } from '../lite-controller';

export interface AutocompleteFieldFormComponentProps {
    name: string;
    label: string;
    placeholder?: string;
    disabled?: boolean;
    errors?: string;

    search: (search: string) => any;
    options: Options;
    customHeader: React.ReactNode;
    children?: React.ReactNode;

    value: any;
    onChange: <T=string, TT=never>(event: {
        reason: AutocompleteChangeReason;
        option: Option<T, TT> | null;
    }) => void;
    onBlur: (event: any) => void;

    multiple?: boolean;
}
export function AutocompleteBaseFieldFormComponent(props: AutocompleteFieldFormComponentProps) {
    const { name, label, placeholder, disabled, errors, customHeader, search, options, value, onChange, onBlur, children } = props;

    const debounceSearch = useCallback(debounce((searchValue: string) => {
        search(searchValue);
    }, 300), [search]);


    const StyledPopper = styled(Popper)(({ theme }) => ({
        width: '100%',
        zIndex: theme.zIndex.modal,
        borderColor: '#FF0000',//theme.palette.neutral![600],
        '& .MuiAutocomplete-paper': {
            margin: theme.spacing(0),
            borderRadius: 0,
            borderWidth: 0
        },
        '>.MuiPaper-root': {
            borderColor: theme.palette.neutral![600],
            borderWidth: 1,
            borderStyle: 'solid',
            boxShadow: theme.shadows[10]
        }
    }));

    return (
        <Autocomplete
            multiple={props.multiple ?? false}
            id={name}
            fullWidth
            limitTags={2}
            sx={{ width: '100%' }}
            filterOptions={(x) => x}
            options={options}
            autoComplete={false}
            includeInputInList
            disabled={disabled}
            filterSelectedOptions={false}
            isOptionEqualToValue={(option, value) => option.value === value.value}
            value={value ?? null}
            noOptionsText="No data available"
            onChange={(event: any, newValue: Option | null, reason: AutocompleteChangeReason, details: any) => {
                // setOptions(newValue ? [newValue, ...options] : options);
                // setValue(newValue);
                
                onChange({
                    reason: reason,
                    option: details?.option
                });
                //search(newValue);
            }}
            onInputChange={(event, newInputValue) => {
                //setInputValue(newInputValue);
                debounceSearch(newInputValue);
            }}
            renderInput={(params) => (
                <TextField {...params}
                    label={label}
                    fullWidth sx={{ width: '100%' }}
                    error={!!errors}
                    helperText={errors} />
            )}
            slots={{
                popper: (popperProps: PopperProps & AutocompletePopperSlotPropsOverrides) => {

                    return <StyledPopper {...popperProps} onMouseDown={event => event.preventDefault()} placement="bottom-start">
                        <Paper>

                            {customHeader}

                            {typeof popperProps.children === 'function'
                                ? popperProps.children({ placement: popperProps.placement! })
                                : popperProps.children}

                        </Paper>
                    </StyledPopper>
                }
            }}

            renderOption={(props, option: Option) => {
                const { key, ...optionProps } = props;

                return (
                    <Box
                        key={JSON.stringify(option.value)}
                        component="li"
                        sx={{ '& > img': { mr: 2, flexShrink: 0 } }}
                        {...optionProps}
                    >

                        <Typography variant="body2" color="text.secondary">
                            {option.label}
                        </Typography>
                    </Box>
                );
            }}

        // renderTags={(_, selectedOptions) => 'test'}
        />
    );
}




// export type AutocompleteOptionValue = PrimitiveOption<any>;