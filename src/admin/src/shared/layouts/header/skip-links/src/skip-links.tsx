import { makeVar, useReactiveVar } from "@apollo/client";
import { Visibility } from "@mui/icons-material";
import { Box, BoxProps, Button, styled } from "@mui/material";
import { useEffect } from "react";


const LinkStyle = styled(Button)(({ theme }) => ({
    backgroundColor: theme.palette.background.paper,
    borderRadius: 1,

    position: 'absolute',
    left: '-10000px',
    width: '1px',
    height: '1px',
    margin: '0px 10px',
    overflow: 'hidden',
    zIndex: 1,
    
    '&:hover': {
        backgroundColor: theme.palette.background.paper,
    },

    '&:focus': {
        left: '0',
        width: 'auto',
        height: 'auto',
        overflow: 'visible',
    },

    '&:active': {
        left: '0',
        width: 'auto',
        height: 'auto',
        overflow: 'visible',
    },



}));

const SkipLinkContainer = styled(Box)<BoxProps>(({ theme }) => ({

    ul: {
        listStyle: 'none',
        margin: 0,
        padding: 0,
    }
})) as typeof Box;

interface Link {
    title: string;
    id: string;
}


const skipLinksVar = makeVar<Link[]>([{
    title: 'Skip Navigation',
    id: 'main'
}]);

export function SkipLinks() {
    const links = useReactiveVar(skipLinksVar);

    return (
        <SkipLinkContainer component="nav">
            <ul>
                {links.map((link: Link) => (
                    <li key={link.title}>
                        <LinkStyle variant="outlined" href={`#${link.id}`}>
                            {link.title}
                        </LinkStyle>
                    </li>
                ))}
            </ul>
        </SkipLinkContainer>
    )
}

interface SkipLinkProps {
    id: string;
    title: string;
}
export function useSkipLink(props: SkipLinkProps) {

    useEffect(() => {
        skipLinksVar([...skipLinksVar(), { id: props.id, title: props.title }]);

        return () => {
            skipLinksVar(skipLinksVar().filter(i => i.id != props.id));
        }
    }, []);

    return props.id;
}