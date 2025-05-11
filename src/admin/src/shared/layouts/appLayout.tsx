"use client";

import { Header, loggedInHeaderProps } from "@/shared/layouts/header/Header";
import { useState } from "react";
import { Topnav } from "./header/Topnav/top-nav";
import { usePinNav } from "@/lib/pinNav";
import { Container, ContentLayout, LeftSidePanel } from "./rootLayout";


interface Props {
    children: React.ReactNode;
    sideChildren?: React.ReactNode;
}
export default function AppLayout(props: Props) {

    const [mobileNavOpen, setMobileNavOpen] = useState(false);

    return (<Container>
        <Header
            nav={<Topnav />}
            rightItems={loggedInHeaderProps.rightItems}
            hasLeftPanel={!!props.sideChildren}
            onNavOpen={() => setMobileNavOpen(true)}
        />

        <LeftSidePanel
            mobileNavOpen={mobileNavOpen}
            setMobileNavOpen={setMobileNavOpen}
            mobileNav={<Topnav direction="column" />}>

            {props.sideChildren}
        </LeftSidePanel>

        <ContentLayout id="main">
            {props.children}
        </ContentLayout>

    </Container>)
}
