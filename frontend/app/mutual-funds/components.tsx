"use client";

import * as React from "react";
import { Check, ChevronsUpDown } from "lucide-react";

import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import {
	Command,
	CommandEmpty,
	CommandGroup,
	CommandInput,
	CommandItem,
	CommandList,
} from "@/components/ui/command";
import {
	Popover,
	PopoverContent,
	PopoverTrigger,
} from "@/components/ui/popover";
import { Scheme } from "./types";

export function ComboboxDemo({
	setSchemeId,
}: {
	setSchemeId: (schemeId: number) => void;
}) {
	const [open, setOpen] = React.useState(false);
	const [selectedScheme, setSelectedScheme] = React.useState("");

	const [availableSchemes, setAvailableSchemes] = React.useState<Scheme[]>(
		[]
	);
	const [schemes, setSchemes] = React.useState<Scheme[]>([]);
	const [searchText, setSearchText] = React.useState("");

	React.useEffect(() => {
		// send http request to /mf-schemes.json
		fetch("/mf-schemes.json")
			.then((res) => res.json())
			.then((data) => {
				setAvailableSchemes(data);
				setSchemes(data.slice(0, 10));
			});
	}, []);

	React.useEffect(() => {
		const filtered = availableSchemes
			.filter((scheme) => {
				return scheme.SchemeName.toLowerCase().includes(
					searchText.toLowerCase()
				);
			})
			.slice(0, 10);

		setSchemes(filtered);
	}, [searchText]);

	return (
		<div className="pr-4 w-full">
			<Popover open={open} onOpenChange={setOpen}>
				<PopoverTrigger asChild>
					<Button
						variant="outline"
						role="combobox"
						aria-expanded={open}
					>
						{selectedScheme
							? schemes.find(
									(scheme) =>
										scheme.SchemeName === selectedScheme
							  )?.SchemeName
							: "Select mutual fund scheme..."}
						<ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
					</Button>
				</PopoverTrigger>
				<PopoverContent className="p-0">
					<Command>
						<CommandInput
							placeholder="Search..."
							onChangeCapture={(e) =>
								setSearchText(e.currentTarget.value)
							}
						/>
						<CommandList>
							<CommandEmpty>No scheme found.</CommandEmpty>
							<CommandGroup>
								{schemes.map((scheme) => (
									<CommandItem
										key={scheme.ID}
										value={scheme.SchemeName}
										onSelect={(currentValue) => {
											setSelectedScheme(
												currentValue === selectedScheme
													? ""
													: currentValue
											);
											setSchemeId(scheme.ID);
											setOpen(false);
										}}
									>
										<Check
											className={cn(
												"mr-2 h-4 w-4",
												selectedScheme ===
													scheme.SchemeName
													? "opacity-100"
													: "opacity-0"
											)}
										/>
										{scheme.SchemeName}
									</CommandItem>
								))}
							</CommandGroup>
						</CommandList>
					</Command>
				</PopoverContent>
			</Popover>
		</div>
	);
}
