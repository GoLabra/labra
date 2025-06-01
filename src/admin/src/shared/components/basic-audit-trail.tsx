import { useMemo } from "react";
import { HistoryEvent, HistoryTimeLine } from "./history-timeline";
import { Stack, Typography } from "@mui/material";
import { localeConfig } from "@/config/locale-config";
import { InlineAvatar } from "./avatar";
import dayjs from "dayjs";

interface BasicAuditTrailProps {
	defaultValues: any;
}
export const BasicAuditTrail = ({ defaultValues }: BasicAuditTrailProps) => {

	if (!defaultValues) {
		defaultValues = {};
	}

	const { createdBy, createdAt, updatedBy, updatedAt } = defaultValues;
	const timelineEvents = useMemo((): HistoryEvent[] => {
		return [{
			label: 'Last Updated',
			details: (
				<Stack>
					<Typography
						color="text.secondary"
						fontSize={12}
						lineHeight='1.5rem'>
						{dayjs(updatedAt).format(localeConfig.dateTime.displayFormat)}
					</Typography>

					<Stack direction="row" gap={1}>
						<InlineAvatar name={updatedBy?.name ?? 'Unknown'} />
						<Typography color="text.secondary">{updatedBy?.name ?? 'Unknown'}</Typography>
					</Stack>

				</Stack>)
		}, {
			label: 'Created',
			details: (
				<Stack>
					<Typography
						color="text.secondary"
						fontSize={12}
						lineHeight='1.5rem'>
						{dayjs(createdAt).format(localeConfig.dateTime.displayFormat)}
					</Typography>

					<Stack direction="row" gap={1}>
						<InlineAvatar name={createdBy?.name ?? 'Unknown'} />
						<Typography color="text.secondary">{createdBy?.name ?? 'Unknown'}</Typography>
					</Stack>
					
				</Stack>)

		}];
	}, [createdBy, createdAt, updatedBy, updatedAt]);

	if (!defaultValues) {
		return null;
	}

	return (
		<HistoryTimeLine events={timelineEvents} />
	)
}