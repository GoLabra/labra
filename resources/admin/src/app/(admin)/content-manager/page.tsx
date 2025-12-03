"use client";

import { useDocumentTitle } from "@/hooks/use-document-title";
import { useCurrentEntityNameContext } from "@/hooks/use-current-entity";
import { ContentManagerHome } from "@/features/content-manager/content-manager.home";
import ContentManager from "@/features/content-manager/content-manager";


export default function Page() {
	useDocumentTitle({ title: "Content Manager" });

	const entityName = useCurrentEntityNameContext();

	if(entityName) {
		return <ContentManager />
	}

	return <ContentManagerHome />

}
