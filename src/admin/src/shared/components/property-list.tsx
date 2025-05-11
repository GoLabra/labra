"use client";

import type { FC, ReactNode } from 'react';
import PropTypes from 'prop-types';
import { List } from '@mui/material';

interface PropertyListProps {
  children: ReactNode;
}

export const PropertyList: FC<PropertyListProps> = (props) => {
  const { children } = props;

  return (
    <List disablePadding>
      {children}
    </List>
  );
};
