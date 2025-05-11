"use client";

import type { ChangeEvent, FC, FocusEvent, KeyboardEvent } from 'react';
import { useCallback, useEffect, useRef, useState } from 'react';
import PropTypes from 'prop-types';
import MagnifyingGlassIcon from '@heroicons/react/24/outline/MagnifyingGlassIcon';
import type { SxProps } from '@mui/system';
import { InputBase, SvgIcon } from '@mui/material';
import { styled } from '@mui/material/styles';

const QueryFieldRoot = styled('div')(
  (({ theme }) => ({
    alignItems: 'center',
    backgroundColor: 'background.paper',
    border: `1px solid ${theme.palette.divider}`,
    borderRadius: theme.shape.borderRadius,
    display: 'flex',
    height: 42,
    padding: '0 16px'
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
  const [autoFocus, setAutoFocus] = useState<boolean>(false);
  const inputRef = useRef<HTMLInputElement | null>(null);
  const [value, setValue] = useState<string>('');

  useEffect(
    () => {
      setValue(initialValue);
    },
    [initialValue]
  );

  useEffect(
    () => {
      if (!disabled && autoFocus && inputRef?.current) {
        inputRef.current.focus();
      }
    },
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [disabled]
  );

  const handleChange = useCallback(
    (event: ChangeEvent<HTMLInputElement>): void => {
      setValue(event.target.value);
    },
    []
  );

  const handleKeyup = useCallback(
    (event: KeyboardEvent<HTMLInputElement>): void => {
      if (event.code === 'Enter') {
        onChange?.(value);
      }
    },
    [value, onChange]
  );

  const handleFocus = useCallback(
    (): void => {
      setAutoFocus(true);
    },
    []
  );

  const handleBlur = useCallback(
    (event: FocusEvent<HTMLInputElement>): void => {
      /*
       There is a situation where an input goes from not disabled to disabled and DOM emits a blur
       event, with event as undefined. This means, that sometimes we'll receive an React Synthetic
       event and sometimes undefined because when DOM triggers the event, React is unaware of it,
       or it simply does not emit the event. To bypass this behaviour, we store a local variable
       that acts as autofocus.
       */

      if (event) {
        setAutoFocus(false);
      }
    },
    []
  );

  return (
    <QueryFieldRoot {...other}>
      <SvgIcon
        fontSize="small"
        sx={{ mr: 1 }}
      >
        <MagnifyingGlassIcon />
      </SvgIcon>
      <InputBase
        disabled={disabled}
        inputProps={{ ref: inputRef,
            'aria-label': "Search grid"
         }}
        onBlur={handleBlur}
        onChange={handleChange}
        onFocus={handleFocus}
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
