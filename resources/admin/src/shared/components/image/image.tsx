import { forwardRef } from 'react';
import { LazyLoadImage } from 'react-lazy-load-image-component';

import { Box, BoxProps } from '@mui/material';
import { styled } from '@mui/material/styles';

// import { CONFIG } from 'src/config-global';

import { imageClasses } from './classes';

// ----------------------------------------------------------------------

const ImageWrapper = styled(Box)({
    overflow: 'hidden',
    position: 'relative',
    verticalAlign: 'bottom',
    display: 'inline-block',
    [`& .${imageClasses.wrapper}`]: {
        width: '100%',
        height: '100%',
        verticalAlign: 'bottom',
        backgroundSize: 'cover !important',
    },
}) as typeof Box;

const Overlay = styled('span')({
    top: 0,
    left: 0,
    zIndex: 1,
    width: '100%',
    height: '100%',
    position: 'absolute',
});

// ----------------------------------------------------------------------

interface ImageProps extends BoxProps {
    ratio?: number;
    disabledEffect?: boolean;
    //
    alt: string;
    src: string;
    delayTime?: number;
    threshold?: number;
    beforeLoad?: () => void;
    delayMethod?: 'throttle' | 'debounce';
    placeholder?: string;
    wrapperProps?: any;
    scrollPosition?: 'top' | 'center' | 'bottom';
    effect?: 'blur' | 'grayscale' | 'sepia' | 'hue-rotate' | 'invert' | 'opacity';
    visibleByDefault?: boolean;
    wrapperClassName?: string;
    useIntersectionObserver?: boolean;
    //
    slotProps?: {
        overlay?: any;
    };
    sx?: any;
}
export const Image = forwardRef(
    (
        {
            ratio,
            disabledEffect = false,
            //
            alt,
            src,
            delayTime,
            threshold,
            beforeLoad,
            delayMethod,
            placeholder,
            wrapperProps,
            scrollPosition,
            effect = 'blur',
            visibleByDefault,
            wrapperClassName,
            useIntersectionObserver,
            //
            slotProps,
            sx,
            ...other
        }: ImageProps,
        ref
    ) => {
        const content = (
            <Box
                component={LazyLoadImage as any}
                alt={alt}
                src={src}
                delayTime={delayTime}
                threshold={threshold}
                beforeLoad={beforeLoad}
                delayMethod={delayMethod}
                placeholder={placeholder}
                wrapperProps={wrapperProps}
                scrollPosition={scrollPosition}
                visibleByDefault={visibleByDefault}
                effect={visibleByDefault || disabledEffect ? undefined : effect}
                useIntersectionObserver={useIntersectionObserver}
                wrapperClassName={wrapperClassName || imageClasses.wrapper}
                // placeholderSrc={
                //   visibleByDefault || disabledEffect
                //     ? `${CONFIG.site.basePath}/assets/transparent.png`
                //     : `${CONFIG.site.basePath}/assets/placeholder.svg`
                // }
                sx={{
                    width: 1,
                    height: 1,
                    objectFit: 'cover',
                    verticalAlign: 'bottom',
                    aspectRatio: ratio,
                }}
            />
        );

        return (
            <ImageWrapper
                ref={ref}
                component="span"
                className={imageClasses.root}
                sx={{ ...(!!ratio && { width: 1 }), ...sx }}
                {...other}
            >
                {slotProps?.overlay && <Overlay className={imageClasses.overlay} sx={slotProps?.overlay} />}

                {content}
            </ImageWrapper>
        );
    }
);
Image.displayName = 'Image';