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
import { AddInvestmentResponse, SchemeData } from "./types";
import { useToast } from "@/hooks/use-toast";

export default function MutualFunds() {
	const [schemeId, setSchemeId] = React.useState(-1);
	const [nav, setNav] = React.useState<number>();
	const [units, setUnits] = React.useState<number>();
	const [investedAt, setInvestedAt] = React.useState<Date>();

	const { toast } = useToast();

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
						onChange={(e) => setNav(Number(e.target.value))}
					/>
					<Input
						className="mr-3 w-30"
						type="text"
						placeholder="Units"
						onChange={(e) => setUnits(Number(e.target.value))}
					/>
					<DatePicker date={investedAt} setDate={setInvestedAt} />
				</div>
				<div className="flex flex-row mt-5">
					<Button
						onClick={() => {
							fetch("/backend/api/mf/investment/add", {
								method: "POST",
								headers: {
									"Content-Type": "application/json",
								},
								body: JSON.stringify({
									scheme_id: schemeId,
									nav: nav,
									units: units,
									invested_at: investedAt,
								}),
							})
								.then((response) => response.json())
								.then((response: AddInvestmentResponse) => {
									if (response.message == "success") {
										toast({
											description: "Investment added",
										});
									} else {
										toast({
											description: "Investment added",
											variant: "destructive",
										});
									}
								});
						}}
					>
						Add new investment
					</Button>
				</div>
			</div>
		</div>
	);
}

function DatePicker({
	date,
	setDate,
}: {
	date: Date | undefined;
	setDate: React.Dispatch<React.SetStateAction<Date | undefined>>;
}) {
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
	const [schemeData, setSchemeData] = React.useState<SchemeData>();

	React.useEffect(() => {
		if (schemeId === -1) {
			return;
		}
		// send http request to get nav data for scheme id
		fetch("/backend/api/mf/schemes/" + schemeId)
			.then((res) => res.json())
			.then((data: SchemeData) => {
				setSchemeData(data);
			});
	}, [schemeId]);

	if (schemeId === -1) {
		return <div></div>;
	}

	return (
		<div className="flex flex-row">
			<div className="w-50">
				Current nav value: {schemeData?.currentNav}
			</div>
			<div>Date: {schemeData?.date}</div>
		</div>
	);
}
