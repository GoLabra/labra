"use client";

import type { FC } from 'react';
import PropTypes from 'prop-types';
import PlusIcon from '@heroicons/react/24/outline/PlusIcon';
import type { SxProps } from '@mui/system';
import { Box, Button, SvgIcon, Typography } from '@mui/material';
import { styled } from '@mui/material/styles';
import Image from 'next/image';

const ResourceUnavailableRoot = styled('div')(
  ({ theme }) => ({
    alignItems: 'center',
    // backgroundColor: theme.palette.mode === 'dark'
    //   ? theme.palette.info[900]
    //   : theme.palette.info[50],
    display: 'flex',
    flexDirection: 'column',
    justifyContent: 'center',
    padding: theme.spacing(3)
  })
);

interface ResourceUnavailableProps {
  message?: string;
  onCreate?: () => void;
  sx?: SxProps;
}

export const ResourceUnavailable: FC<ResourceUnavailableProps> = (props) => {
  const { message = '', onCreate, sx } = props;

  return (
    <ResourceUnavailableRoot sx={sx}>
      <Box
        sx={{
          '& img': {
            maxWidth: '100%'
          }
        }}
      >
        <Image src="/assets/illustration-not-found.svg" alt="Illustration not found" />
      </Box>
      {message && (
        <Typography
          color="text.secondary"
          sx={{ mt: 2 }}
          variant="body2"
        >
          {message}
        </Typography>
      )}
      {onCreate && (
        <Button
          onClick={onCreate}
          startIcon={(
            <SvgIcon fontSize="small">
              <PlusIcon />
            </SvgIcon>
          )}
          sx={{ mt: 2 }}
          variant="contained"
        >
          Create Resource
        </Button>
      )}
    </ResourceUnavailableRoot>
  );
};

ResourceUnavailable.propTypes = {
  message: PropTypes.string,
  onCreate: PropTypes.func,
  // @ts-ignore
  sx: PropTypes.object
};
