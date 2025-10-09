import { Accordion, AccordionDetails, AccordionSummary, DialogContent, Stack, Typography } from "@mui/material"
import React, { ReactElement, useCallback, useEffect } from "react";
import { ReactNode, SyntheticEvent, useState } from "react";
import { Scrollbar } from "./scrollbar";
import { useMounted } from "@/hooks/use-mounted";
import useDidMountEffect from "@/hooks/use-did-mount-effect ";


interface AccordeonStepperProps {
    autoselectNextStepWhenCompleted: true;
    children: React.ReactElement<AccordeonStepperStepProps>[];
}
AccordeonStepper.defaultProps = {
    autoselectNextStepWhenCompleted: true // Default value
};
export function AccordeonStepper(props: AccordeonStepperProps) {

    const { autoselectNextStepWhenCompleted, children } = props;

    const [selectedIndex, setSelectedIndex] = useState<number>(0);


    const onStepSelected = useCallback(
        (index: number) => {

            //check prev step if completed
            if(index == 0){
                setSelectedIndex(index);
                return;
            }
            
            if( children[index-1].props.completed) {
                setSelectedIndex(index);
                return;
            }
        },
        [setSelectedIndex, children]
    )

    const onCompletedChanged = useCallback(
        (index: number) => {

            if(index != selectedIndex){
                return;
            }

            if(autoselectNextStepWhenCompleted){
                setSelectedIndex((prev:number) => {
                    return index+1;
                })
            }
        },
        [autoselectNextStepWhenCompleted, selectedIndex, setSelectedIndex]
    )
    
    const childrenWithProps = React.Children.map(children, (child, index) => {
        return React.cloneElement(child, {
                                                onSelected: () => onStepSelected(index),
                                                selected: selectedIndex === index,
                                                indexDisplay: index + 1,
                                                onCompletedChanged: () => onCompletedChanged(index)
                                            }
        );
    });

    return (
        <>
            {childrenWithProps}
        </>
    )
}

interface AccordeonStepperStepProps {
    label: string,
    onSelected?: () => void;
    onCompletedChanged?: () => void;
    selected?: boolean;
    indexDisplay?: ReactNode;
    completed?: boolean;
    completedMark?: ReactNode;
    children?: ReactNode;
}
export const AccordeonStepperStep = (props: AccordeonStepperStepProps) => {

    const { label, onSelected, selected, indexDisplay, completed, onCompletedChanged, completedMark, children } = props;

    useEffect(() => {
        if(completed){
            onCompletedChanged?.();
        }
    }, [completed, onCompletedChanged])

    return (
        <Accordion expanded={selected} onChange={onSelected}>

            <AccordionSummary aria-controls="panel1d-content" id="panel1d-header">

                <Stack
                    direction="row"
                    alignItems="center"
                    gap={2}>

                    {indexDisplay}

                    <Stack>
                        <Typography>{label}</Typography>
                        { completed && (completedMark) }
                    </Stack>

                </Stack>

            </AccordionSummary>
            <AccordionDetails>
                
                {children}
                
            </AccordionDetails>


        </Accordion>
    )
}