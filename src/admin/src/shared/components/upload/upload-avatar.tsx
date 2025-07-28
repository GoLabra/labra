import { useEffect, useState } from 'react';
import { useDropzone } from 'react-dropzone';
import { Image } from "@/shared/components/image/image";
import { Box, Typography } from '@mui/material';
import CameraAltIcon from '@mui/icons-material/CameraAlt';
import { AvatarShell } from '../avatar';

interface UploadAvatarProps {
    disabled?: boolean;
    value: any;
    errors?: any;
    helperText?: string;
    onDrop?: (acceptedFiles: any) => void;
}
export const UploadAvatar = (props: UploadAvatarProps) => {


    const dropzone = useDropzone({
        multiple: false,
        disabled: props.disabled,
        accept: { 'image/*': [] },
        onDrop: props.onDrop
    });

    const hasFile = !!props.value;
    const hasError = dropzone.isDragReject || !!props.errors;

    const [preview, setPreview] = useState('');

    const renderPreview = hasFile && (
        <Image alt="avatar" src={preview} sx={{ width: 1, height: 1, borderRadius: '100%' }} />
    );

    useEffect(() => {
        if (typeof props.value === 'string') {
            setPreview(props.value);
        } else if (props.value instanceof File) {
            setPreview(URL.createObjectURL(props.value));
        }
    }, [props.value]);

    const renderPlaceholder = (
        <Box
            className="upload-placeholder"
            sx={{
                top: 0,
                gap: 1,
                left: 0,
                width: 1,
                height: 1,
                zIndex: 9,
                display: 'flex',
                borderRadius: '50%',
                position: 'absolute',
                alignItems: 'center',
                color: 'text.disabled',
                flexDirection: 'column',
                justifyContent: 'center',
                bgcolor: 'color-mix(in srgb, var(--mui-palette-grey-500), transparent 80%)',
                transition: (theme) =>
                    theme.transitions.create(['opacity'], { duration: theme.transitions.duration.shorter }),
                '&:hover': { opacity: 0.72 },
                ...(hasError && {
                    color: 'error.main',
                     bgcolor: 'color-mix(in srgb, var(--mui-palette-error-main), transparent 80%)',
                }),
                ...(hasFile && {
                    zIndex: 9,
                    opacity: 0,
                    color: 'common.white',
                    bgcolor: 'color-mix(in srgb, var(--mui-palette-grey-900), transparent 50%)',
                }),
            }}
        >
            <CameraAltIcon />
            <Typography variant="caption">{hasFile ? 'Change photo' : 'Upload photo'}</Typography>
        </Box>
    );

    const renderContent = (
        <Box
            sx={{
                width: 1,
                height: 1,
                overflow: 'hidden',
                borderRadius: '50%',
                position: 'relative',
            }}
        >
            {renderPreview}
            {renderPlaceholder}
        </Box>
    );

    return (
        <>
            <AvatarShell
                {...dropzone.getRootProps()}
                sx={{
                    cursor: 'pointer',
                    ...(dropzone.isDragActive && { opacity: 0.72 }),
                    ...(props.disabled && { opacity: 0.48, pointerEvents: 'none' }),
                    ...(hasError && { borderColor: 'error.main' }),
                    ...(hasFile && {
                        ...(hasError && {
                            bgcolor: 'color-mix(in srgb, var(--mui-palette-error-main), transparent 80%)',
                        }),
                        '&:hover .upload-placeholder': { opacity: 1 },
                    }),
                    //...props.sx,
                }}
            >
                <input {...dropzone.getInputProps()} />

                {renderContent}
            </AvatarShell>

            {props.helperText && props.helperText}

            {/* <RejectionFiles files={dropzone.fileRejections} /> */}
        </>
    );
}
