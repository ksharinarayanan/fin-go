"use client";

import { Button } from "@/components/ui/button";
import { ComboboxDemo } from "./components";
import { Input } from "@/components/ui/input";
import {
	Popover,
	PopoverContent,
	PopoverTrigger,
} from "@/components/ui/popover";
import { Calendar } from "@/components/ui/calendar";
import { cn } from "@/lib/utils";
import { Calendar as CalendarIcon } from "lucide-react";
import { format } from "date-fns";
import React from "react";

export default function MutualFunds() {
	const [schemeId, setSchemeId] = React.useState(-1);

	return (
		<div className="flex w-full">
			<div className="flex flex-col p-10">
				<div className="flex flex-row">
					<ComboboxDemo setSchemeId={setSchemeId} />
					<CurrentNav schemeId={schemeId} />
				</div>
				<div className="flex flex-row mt-5">
					<Input
						className="mr-3 w-50"
						type="text"
						placeholder="NAV"
					/>
					<Input
						className="mr-3 w-30"
						type="text"
						placeholder="Units"
					/>
					<DatePickerDemo />
				</div>
				<div className="flex flex-row mt-5">
					<Button>Add new investment</Button>
				</div>
			</div>
		</div>
	);
}

function DatePickerDemo() {
	const [date, setDate] = React.useState<Date>();

	return (
		<Popover>
			<PopoverTrigger asChild>
				<Button
					variant={"outline"}
					className={cn(
						"w-[280px] justify-start text-left font-normal",
						!date && "text-muted-foreground"
					)}
				>
					<CalendarIcon className="mr-2 h-4 w-4" />
					{date ? format(date, "PPP") : <span>Date invested</span>}
				</Button>
			</PopoverTrigger>
			<PopoverContent className="w-auto p-0">
				<Calendar
					mode="single"
					selected={date}
					onSelect={setDate}
					initialFocus
				/>
			</PopoverContent>
		</Popover>
	);
}

function CurrentNav({ schemeId }: { schemeId: number }) {
	const [schemeData, setSchemeData] = React.useState(0);

	React.useEffect(() => {
		if (schemeId === -1) {
			return;
		}
		// send http request to get nav data for scheme id
		fetch("http://localhost:8080/api/mf/schemes/" + schemeId)
			.then((res) => res.json())
			.then((data) => {
				setSchemeData(data);
			});
	}, [schemeId]);

	if (schemeId === -1) {
		return <div></div>;
	}

	return (
		<div className="flex flex-row">
			<div className="w-50">
				Current nav value: {schemeData.currentNav}
			</div>
			<div>Date: {schemeData.date}</div>
		</div>
	);
}
