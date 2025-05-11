import { useTabs } from "@/hooks/useTabs";
import { Tab, Tabs } from "@mui/material"
import { AccountGeneralTab } from "./account-general-tab";
import AssignmentIndIcon from '@mui/icons-material/AssignmentInd';
import SecurityIcon from '@mui/icons-material/Security';

export const UserAccount = () => {

    const tabs = useTabs('general');

    return (<>
        <Tabs {...tabs} sx={{ mb: { xs: 3, md: 5 } }}>
            <Tab label="General" icon={<AssignmentIndIcon />} value="general" />
            <Tab label="Security" icon={<SecurityIcon />} value="security" />
        </Tabs>

        {tabs.value === 'general' && <AccountGeneralTab />}
    </>)
}