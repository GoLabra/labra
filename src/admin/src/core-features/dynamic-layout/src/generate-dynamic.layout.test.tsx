import "@testing-library/jest-dom";
import { render, screen } from "@testing-library/react";
import { Box } from "@mui/material";
import { renderDynamicLayout } from "./generate-dynamic-layout";
import { dynamicLayoutItem } from "./dynamic-layout";


it("renders the dynamic layout", async () => {

    const dynamicLayoutItems: dynamicLayoutItem<string> = {
        type: "group",
        direction: "vertical",
        gap: 1.5,
        key: "001",
        children: [
            {
                data: 'item 01',
            }, {
                data: 'item 02',
            }, {
                type: "group",
                direction: "horizontal",
                key: "002",
                children: [{
                        data: 'item 03',
                    }, {
                        data: 'item 04',
                    }
                ]
            },
            {
                data: 'item 05',
            }
        ],
    };

    const { container } = render(
        <>
            {renderDynamicLayout(dynamicLayoutItems, (data) => <Box key={data}>{data}</Box>)}
        </>
    );

    const mainGroup = container.querySelector('div > div');
    const numberOfChildren = mainGroup!.childElementCount;
    expect(numberOfChildren).toBe(4);

    const horizontalGroup = mainGroup!.querySelector('div:nth-child(3)');
    const numberOfChildren2 = horizontalGroup!.childElementCount;
    expect(numberOfChildren2).toBe(2);

    for(let data of ['item 01', 'item 02', 'item 03', 'item 04', 'item 05']){
        const dataItem = screen.getByText(data);
        expect(dataItem).toBeInTheDocument();
    }
});