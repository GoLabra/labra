'use client'

import XCircleIcon from '@heroicons/react/24/outline/XCircleIcon';
import {
    alpha,
    createTheme,
    filledInputClasses,
    inputAdornmentClasses,
    inputBaseClasses,
    inputLabelClasses,
    PaletteColor,
    SvgIcon,
    switchClasses,
    tableCellClasses,
    Theme
} from '@mui/material';
import type { Components } from '@mui/material/styles/components';

// Used only to create transitions
const muiTheme = createTheme();


export const bgCyanBlur = {
    backgroundImage: 'url(/assets/core/cyan-blur.png)',
    backgroundRepeat: 'no-repeat',
    backgroundSize: '100%',
    backgroundPosition: 'right top',
};

export const bgRedBlur = {
    backgroundImage: 'url(/assets/core/red-blur.png)',
    backgroundRepeat: 'no-repeat',
    backgroundSize: '50%',
    backgroundPosition: 'left bottom',
}

export const bgCyanRedBlur = {
    backgroundImage: 'url(/assets/core/cyan-blur.png), url(/assets/core/red-blur.png)',
    backgroundRepeat: 'no-repeat, no-repeat',
    backgroundSize: '50%, 50%',
    backgroundPosition: 'right top, left bottom',
}


export const createComponents = (): Components => {
    return {
        MuiAppBar: {
            styleOverrides: {
                root: {
                    backgroundColor: 'var(--mui-palette-background-paper)'
                }
            }
        },
        MuiAutocomplete: {
            styleOverrides: {
                root: {
                    [`& .${filledInputClasses.root}`]: {
                        paddingTop: 6
                    },
                    '.MuiFilledInput-root input.MuiFilledInput-input': {
                        padding: 0,
                        height: 'unset',
                        fontSize: 14,
                        fontWeight: 500,
                        lineHeight: 1.6
                    },

                    noOptions: {
                        fontSize: 14,
                        letterSpacing: 0.15,
                        lineHeight: 1.6
                    },
                    option: {
                        fontSize: 14,
                        letterSpacing: 0.15,
                        lineHeight: 1.6
                    }
                }
            }
        },
        MuiAvatar: {
            styleOverrides: {
                root: {
                    fontSize: 14,
                    fontWeight: 400,
                    letterSpacing: 0
                }
            }
        },
        // MuiToggleButton: {
        //     styleOverrides: {
        //         root: {
        //             backgroundColor: '#FF0000',
        //         },
        //         selected: {
        //             backgroundColor: '#FF0000',
        //         },
        //         disabled: {

        //         },

        //     }
        // },
        MuiButton: {
            defaultProps: {
                disableRipple: true
            },
            styleOverrides: {
                root: ({ theme }: { theme: any }) => ({
                    fontWeight: 500,
                    borderRadius: 3,
                    textTransform: 'none',
                    '&:focus': {
                        boxShadow: `${alpha((theme.palette.primary as PaletteColor).main, 0.25)} 0 0 0 0.2rem`
                    }
                }),
                sizeLarge: {
                    fontSize: 16
                },
                sizeMedium: {
                    fontSize: 15
                },
                sizeSmall: {
                    fontSize: 14
                }
            }
        },
        MuiButtonGroup: {
            defaultProps: {
                disableRipple: true
            }
        },

        MuiChip: {
            defaultProps: {
                deleteIcon: (
                    <SvgIcon>
                        <XCircleIcon />
                    </SvgIcon >
                )
            },
            styleOverrides: {
                // avatar: {
                //     borderRadius: 4
                // },
                root: {
                    borderRadius: 4,
                    fontWeight: 400,
                    letterSpacing: 0
                }
            }
        },
        MuiCssBaseline: {
            styleOverrides: {
                '*': {
                    boxSizing: 'border-box'
                },
                html: {
                    MozOsxFontSmoothing: 'grayscale',
                    WebkitFontSmoothing: 'antialiased',
                    display: 'flex',
                    flexDirection: 'column',
                    minHeight: '100%',
                    width: '100%'
                },
                body: {
                    display: 'flex',
                    flex: '1 1 auto',
                    flexDirection: 'column',
                    minHeight: '100%',
                    width: '100%'
                },
                '#__next': {
                    display: 'flex',
                    flex: '1 1 auto',
                    flexDirection: 'column',
                    height: '100%',
                    width: '100%'
                },
                '#nprogress': {
                    pointerEvents: 'none'
                },
                '#nprogress .bar': {
                    backgroundColor: '#12B76A',
                    height: 3,
                    left: 0,
                    position: 'fixed',
                    top: 0,
                    width: '100%',
                    zIndex: 2000
                }
            }
        },
        MuiTypography: {
            styleOverrides: {
                h1: { fontSize: '1.5rem', fontWeight: 400 }, // Customize size for h1
                h2: { fontSize: '1.3rem', fontWeight: 400 },   // Customize size for h2
                h3: { fontSize: '1.2rem', fontWeight: 400 },// Customize size for h3
                h4: { fontSize: '1rem' }, // Customize size for h4
                h5: { fontSize: '0.9.rem' },// Customize size for h5
                h6: { fontSize: '0.8rem' },   // Customize size for h6
            }
        },
        // MuiAlert: {
        //     styleOverrides: {
        //         root: {
        //             backgroundColor: 'var(--mui-palette-background-paper)',
        //             color: 'var(--mui-palette-text-primary)',
        //             borderRadius: 3,
        //             boxShadow: '0px 1px 2px 0px rgba(0, 0, 0, 0.24)',
        //         },
        //     }
        // },
        MuiDialogActions: {
            styleOverrides: {
                root: {
                    paddingBottom: 10,
                    paddingLeft: 10,
                    paddingRight: 10,
                    paddingTop: 10,
                    borderBottomLeftRadius: 6,
                    borderBottomRightRadius: 6,

                    backgroundColor: 'var(--mui-palette-background-default)',

                    '&>:not(:first-of-type)': {
                        marginLeft: 16
                    }
                }
            }
        },
        MuiDialogContent: {
            styleOverrides: {
                root: {
                    paddingBottom: 20,
                    paddingLeft: 20,
                    paddingRight: 20,
                    marginTop: 20,

                    backgroundColor: 'var(--mui-palette-background-paper)',

                    '&.MuiDialogContent-root': {
                        marginTop: 0
                    },
                }
            }
        },
        MuiDialogTitle: {
            styleOverrides: {
                root: {
                    fontSize: 18,
                    fontWeight: 400,
                    paddingBottom: 12,
                    paddingLeft: 12,
                    paddingRight: 15,
                    paddingTop: 12,
                    borderTopLeftRadius: 6,
                    borderTopRightRadius: 6,

                    backgroundColor: 'var(--mui-palette-background-default)',
                }
            }
        },
        MuiFormControlLabel: {
            styleOverrides: {
                label: {
                    fontSize: 14,
                    letterSpacing: 0.15,
                    lineHeight: 1.43
                }
            }
        },
        MuiIcon: {
            styleOverrides: {
                fontSizeLarge: {
                    fontSize: 32
                }
            }
        },
        MuiIconButton: {
            styleOverrides: {
                root: {
                    //borderRadius: 6,
                    padding: 8
                },
                sizeSmall: {
                    padding: 4
                }
            }
        },
        MuiInputAdornment: {
            styleOverrides: {
                root: {
                    [`&.${inputAdornmentClasses.positionStart}.${inputAdornmentClasses.filled}`]: {
                        '&:not(.MuiInputAdornment-hiddenLabel)': {
                            marginTop: 0
                        }
                    }
                }
            }
        },
        MuiInputBase: {
            styleOverrides: {
                input: {
                    '&::placeholder': {
                        opacity: 1
                    },
                    [`label[data-shrink=false] + .${inputBaseClasses.formControl} &`]: {
                        '&::placeholder': {
                            opacity: 1 + '!important'
                        }
                    }
                }
            }
        },
        MuiFilledInput: {
            styleOverrides: {
                root: ({ theme }: { theme: any }) => ({
                    borderRadius: 3,
                    borderStyle: 'solid',
                    borderWidth: 1,
                    overflow: 'hidden',
                    padding: '6px 12px',

                    transition: muiTheme.transitions.create([
                        'border-color',
                        'box-shadow'
                    ]),
                    '&:before': {
                        display: 'none'
                    },
                    '&:after': {
                        display: 'none'
                    },

                    backgroundColor: theme.palette.background!.paper,
                    [`&.${filledInputClasses.focused}`]: {
                        backgroundColor: 'transparent',
                        borderColor: (theme.palette.primary as PaletteColor).main,
                        boxShadow: `${alpha((theme.palette.primary as PaletteColor).main, 0.25)} 0 0 0 0.2rem`
                    },
                    [`&.${filledInputClasses.error}`]: {
                        borderColor: (theme.palette.error as PaletteColor).main,
                        boxShadow: `${alpha((theme.palette.error as PaletteColor).main, 0.25)} 0 0 0 0.2rem`
                    },
                    '&:hover': {
                        backgroundColor: theme.palette.action!.hover
                    },

                    ...(theme.palette.mode === 'light' ? {

                        borderColor: theme.palette.neutral![300],
                        boxShadow: `0px 1px 2px 0px ${alpha(theme.palette.neutral![800], 0.08)}`,

                        [`&.${filledInputClasses.disabled}`]: {
                            backgroundColor: theme.palette.action!.disabledBackground,
                            borderColor: theme.palette.neutral![800],
                            boxShadow: 'none'
                        },
                    } : {

                        borderColor: theme.palette.neutral![600],
                        boxShadow: `0px 1px 2px 0px ${alpha(theme.palette.neutral![900], 0.08)}`,

                        [`&.${filledInputClasses.disabled}`]: {
                            backgroundColor: theme.palette.action!.disabledBackground,
                            borderColor: theme.palette.neutral![800],
                            boxShadow: 'none'
                        },
                    }),

                }),
                input: {
                    padding: 0,
                    height: 'unset',
                    fontSize: 14,
                    fontWeight: 500,
                    lineHeight: 1.6
                }
            }
        },
        MuiDrawer: {
            styleOverrides: {
                root: {
                    '>.MuiPaper-root': {
                        backgroundColor: 'var(--mui-palette-background-default)',
                        // borderRightWidth: 0
                    },
                }
            }
        },
        MuiFormLabel: {
            styleOverrides: {
                root: {
                    fontSize: 14,
                    fontWeight: 500,
                    color: 'var(--mui-palette-text-primary)',
                    [`&.${inputLabelClasses.filled}`]: {
                        marginBottom: 8,
                        position: 'relative',
                        transform: 'none'
                    }
                }
            }
        },
        MuiListItemIcon: {
            styleOverrides: {
                root: {
                    minWidth: '36px'
                }
            }
        },
        MuiPaper: {
            styleOverrides: {
                root: {
                    backgroundImage: 'none'
                }
            }
        },
        MuiRadio: {
            styleOverrides: {
                root: {
                    transition: 'color 250ms',
                    '&:hover': {
                        backgroundColor: 'transparent'
                    }
                }
            }
        },
        MuiSelect: {
            defaultProps: {
                variant: 'filled'
            },
            styleOverrides: {
                filled: {
                    '&:focus': {
                        backgroundColor: 'transparent'
                    }
                }
            }
        },
        // MuiSkeleton: {
        //     styleOverrides: {
        //         root: {
        //             borderRadius: 4
        //         }
        //     }
        // },
        MuiSvgIcon: {
            styleOverrides: {
                fontSizeLarge: {
                    fontSize: 32
                }
            }
        },
        MuiSwitch: {
            styleOverrides: {
                root: ({ theme }: { theme: any }) => ({
                    borderRadius: 48,
                    height: 24,
                    marginBottom: 8,
                    marginLeft: 8,
                    marginRight: 8,
                    marginTop: 8,
                    padding: 0,
                    width: 44,

                    [`.${switchClasses.switchBase}`]: {
                        padding: 4,

                        [`&.${switchClasses.checked}`]: {
                            [`&+.${switchClasses.track}`]: {
                                opacity: 1
                            },
                            [`& .${switchClasses.thumb}`]: {
                                opacity: 1
                            }
                        },

                        [`&.${switchClasses.disabled}`]: {
                            [`&+.${switchClasses.track}`]: {
                                backgroundColor: alpha(theme.palette.text!.primary!, 0.08),
                                opacity: 1
                            },
                            [`& .${switchClasses.thumb}`]: {
                                backgroundColor: alpha(theme.palette.text!.secondary!, 0.86),
                                opacity: 1
                            }
                        }
                    },

                    [`.${switchClasses.thumb}`]: {
                        backgroundColor: theme.palette.common.white,
                        height: 16,
                        width: 16
                    },
                    [`.${switchClasses.track}`]: {
                        backgroundColor: theme.palette.neutral![500],
                        opacity: 1
                    },
                    '&:focus-within': {
                        boxShadow: `${alpha((theme.palette.primary as PaletteColor).main, 0.25)} 0 0 0 0.2rem`
                    }
                }),
                // switchBase: {
                //     padding: 4
                // },
                // track: {
                //     //  opacity: 1
                // },
                // thumb: {
                //     height: 16, 
                //     width: 16
                // }
            }
        },
        MuiTabs: {
            styleOverrides: {
                root: {
                    minHeight: 0
                },
                flexContainer: {
                    gap: '20px',
                }
            },
        },
        MuiTab: {
            defaultProps: {
                disableRipple: true,
                iconPosition: 'start'
            },
            styleOverrides: {
                root: {
                    padding: '10px 0',
                    minHeight: 0
                }
            }
        },
        MosaicDataTable: {
            styleOverrides: {
                root: {
                    'table>caption': {
                        display: 'none'
                    },
                    '.MosaicDataTablePaper-root': {
                        boxShadow: 'none'
                    },
                    '.MuiTableContainer-root': {
                        width: '100px',
                        minWidth: '100%',
                    },
                    '.MuiTableCell-head': {
                        padding: 0,
                        height: '100%',

                        '.MuiTableCellDockedDiv-root': {
                            // padding: '10px 10px',
                            // margin: 0,
                            borderTopLeftRadius: 2,
                            borderTopRightRadius: 2,
                        }
                    },
                    '.MuiTableCell-body': {
                        height: '100%'
                    },
                    '.MuiTableCell-root': {
                        padding: 0,

                        // '.MuiTableCellDockedDiv-root': {
                        //     padding: '10px 10px',
                        //     margin: 0
                        // }
                    }
                }
            }
        },
        // MuiTableHead: {
        //     styleOverrides: {
        //         root: {
        //             borderBottomWidth: 1,
        //             borderBottomStyle: 'dashed',

        //             [`&.${tableCellClasses.root}`]: {
        //                 fontSize: 11,
        //                 fontWeight: 600,
        //                 textTransform: 'uppercase'
        //             }
        //         }
        //     }
        // },
        // MuiTableRow: {
        //     styleOverrides: {
        //         root: {
        //             [`&:last-of-type .${tableCellClasses.root}`]: {
        //                 borderWidth: 0
        //             }
        //         }
        //     }
        // },
        // MuiTableCell: {
        //     styleOverrides: {
        //         root: {
        //             padding: '5px',
        //             borderBottomWidth: 1,
        //             borderBottomStyle: 'dashed'
        //         }
        //     }
        // },
        MuiTextField: {
            defaultProps: {
                variant: 'filled'
            }
        }
    };
};
