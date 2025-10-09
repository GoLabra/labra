"use client";

import type { FC, ReactNode } from 'react';
import { alpha, darken } from '@mui/material/styles';
import PropTypes from 'prop-types';
import { List } from '@mui/material';

interface ActionListProps {
  children: ReactNode;
}

export const ActionList: FC<ActionListProps> = (props) => {
  const { children } = props;

  return (
    <List
      dense
    >
      {children}
    </List>
  );
};