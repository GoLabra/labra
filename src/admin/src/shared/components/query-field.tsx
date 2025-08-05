"use client";

import type { ChangeEvent, FC, FocusEvent, KeyboardEvent } from 'react';
import { useCallback, useEffect, useRef, useState } from 'react';
import PropTypes from 'prop-types';
import type { SxProps } from '@mui/system';
import { InputBase, SvgIcon } from '@mui/material';
import { styled } from '@mui/material/styles';
import { Key } from './key-handler/types';
import { Shortcut } from './key-handler/shortcut';

const QueryFieldRoot = styled('div')(
	(({ theme }) => ({
		alignItems: 'center',
		backgroundColor: 'background.paper',
		border: `1px solid var(--mui-palette-divider)`,
		borderRadius: theme.shape.borderRadius,
		display: 'flex',
		height: 42,
		padding: '0 12px'
	}))
);

interface QueryFieldProps {
	disabled?: boolean;
	onChange?: (value: string) => void;
	placeholder?: string;
	sx?: SxProps;
	value?: string;
}

export const QueryField: FC<QueryFieldProps> = (props) => {
	const {
		disabled,
		onChange,
		placeholder,
		value: initialValue = '',
		...other
	} = props;
	
	const inputRef = useRef<HTMLInputElement | null>(null);
	const [value, setValue] = useState<string>('');

	useEffect(() => {
		setValue(initialValue);
	},[initialValue]);

	const handleChange = useCallback((event: ChangeEvent<HTMLInputElement>): void => {
		setValue(event.target.value);
	}, []);

	const handleKeyup = useCallback((event: KeyboardEvent<HTMLInputElement>): void => {
		if (event.code === 'Enter') {
			onChange?.(value);
		}

		if (event.code === 'Escape') {
			if(value.length) {
				setValue('');
				onChange?.('');
			}
			else {
				inputRef.current?.blur();
			}
			event.stopPropagation();
		}
	}, [value, onChange]);

	return (
		<QueryFieldRoot {...other}>
			<Shortcut
				keyToHandle={Key.slash}
				onTriggered={() => inputRef.current?.focus()}
				slotProps={{
					ShortcutViewer: {
						sx: {
							marginRight: '10px',
						}
					}
				}}/>
			<InputBase
				disabled={disabled}
				inputProps={{
					ref: inputRef,
					'aria-label': "Search grid"
				}}
				onChange={handleChange}
				onKeyUp={handleKeyup}
				placeholder={placeholder}
				sx={{ flexGrow: 1 }}
				value={value}
			/>
		</QueryFieldRoot>
	);
};

QueryField.propTypes = {
	disabled: PropTypes.bool,
	onChange: PropTypes.func,
	placeholder: PropTypes.string,
	value: PropTypes.string
};
