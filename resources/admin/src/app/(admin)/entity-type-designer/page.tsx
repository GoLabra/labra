"use client";

import { Box, Container, Stack, Typography } from "@mui/material";
import AddIcon from '@mui/icons-material/Add';
import CreateIcon from '@mui/icons-material/Create';
import { useDocumentTitle } from "@/hooks/use-document-title";
import { OfficialCommunicationChannels } from "@/shared/components/official-communication-channels";
import { Counter, CounterGroup } from "@/shared/components/counter";
import EntityTypeDesignerHome from "@/features/entity-type-designer/entity-type-designer.home";
import { useCurrentEntityNameContext } from "@/hooks/use-current-entity";
import EntityTypeDesigner from "@/features/entity-type-designer/entity-type-designer";


export default function Page() {

	const entityName = useCurrentEntityNameContext();
	
	if(entityName){
		return (<EntityTypeDesigner />);
	}
	
    return <EntityTypeDesignerHome />;
}
