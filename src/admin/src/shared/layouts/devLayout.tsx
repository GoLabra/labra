"use client";

import { Header, loggedInHeaderProps } from "@/shared/layouts/header/Header";
import { useState } from "react";
import { Topnav } from "./header/Topnav/top-nav";
import { Container, ContentLayout, LeftSidePanel } from "./rootLayout";
import { DevelopTopnav } from "./header/DeveloperTopnav/developer-top-nav";


interface Props {
    children: React.ReactNode;
    sideChildren?: React.ReactNode;
}
export default function DevLayout(props: Props) {

    const [mobileNavOpen, setMobileNavOpen] = useState(false);

    return (<Container>
        <Header
            nav={<DevelopTopnav />}
            rightItems={loggedInHeaderProps.rightItems}
            hasLeftPanel={!!props.sideChildren}
            onNavOpen={() => setMobileNavOpen(true)}
        />

        <LeftSidePanel
            mobileNavOpen={mobileNavOpen}
            setMobileNavOpen={setMobileNavOpen}
            mobileNav={<DevelopTopnav direction="column" />}>

            {props.sideChildren}
        </LeftSidePanel>

        <ContentLayout id="main">
            {props.children}
        </ContentLayout>

    </Container>)
}
