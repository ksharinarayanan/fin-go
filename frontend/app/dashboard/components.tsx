import { TrendingDown, TrendingUp } from "lucide-react";
import { MutualFundInvestment } from "./types";
import { Badge } from "@/components/ui/badge";
import {
	AccordionContent,
	AccordionItem,
	AccordionTrigger,
} from "@/components/ui/accordion";
import {
	Table,
	TableBody,
	TableCaption,
	TableCell,
	TableHead,
	TableHeader,
	TableRow,
} from "@/components/ui/table";

const formatCurrency = (value: number) => {
	return new Intl.NumberFormat("en-IN", {
		style: "currency",
		currency: "INR",
		minimumFractionDigits: 2,
		maximumFractionDigits: 2,
	}).format(value);
};

const ProfitLossBadge = ({
	value,
	percentage,
}: {
	value: number;
	percentage: number;
}) => {
	const isPositive = value >= 0;
	return (
		<Badge
			variant="outline"
			className={`ml-2 ${
				isPositive
					? "bg-green-300 text-green-800 border-green-300"
					: "bg-red-300 text-red-800 border-red-500"
			}`}
		>
			{isPositive ? (
				<TrendingUp className="w-3 h-3 mr-1" />
			) : (
				<TrendingDown className="w-3 h-3 mr-1" />
			)}
			{percentage.toFixed(2)}% ({formatCurrency(Math.abs(value))})
		</Badge>
	);
};

// please rename this to a better name
function InvestmentTable({ mf }: { mf: MutualFundInvestment }) {
	return (
		<Table>
			<TableCaption>{mf.scheme_name}</TableCaption>
			<TableHeader>
				<TableRow>
					<TableHead className="font-bold">Units</TableHead>
					<TableHead className="font-bold">Invested On</TableHead>
					<TableHead className="font-bold">Invested NAV</TableHead>
					<TableHead className="font-bold">Invested value</TableHead>
					<TableHead className="font-bold">Current value</TableHead>
					<TableHead className="font-bold text-right">
						Net Profit/Loss
					</TableHead>
				</TableRow>
			</TableHeader>
			<TableBody>
				{mf.investments.map((investment, i) => {
					return (
						<TableRow key={i}>
							<TableCell>{investment.units.toFixed(2)}</TableCell>
							<TableCell>
								{new Intl.DateTimeFormat("en-GB", {
									day: "2-digit",
									month: "short",
									year: "numeric",
								}).format(new Date(investment.invested_at))}
							</TableCell>
							<TableCell>
								{formatCurrency(investment.invested_nav)}
							</TableCell>
							<TableCell>
								{formatCurrency(investment.invested_value)}
							</TableCell>
							<TableCell>
								{formatCurrency(investment.current_value)}
							</TableCell>
							<TableCell className="text-right">
								<ProfitLossBadge
									value={investment.net_profit_loss}
									percentage={
										investment.net_profit_loss_percentage
									}
								/>
							</TableCell>
						</TableRow>
					);
				})}
			</TableBody>
		</Table>
	);
}

export default function MutualFundAccordion({
	mf,
}: {
	mf: MutualFundInvestment;
}) {
	return (
		<AccordionItem value={mf.scheme_id.toString()}>
			<AccordionTrigger className="font-bold">
				<div className="w-full">
					<div className="text-start">{mf.scheme_name}</div>
					<div className="text-end px-4">
						<span className="px-2">
							Day:{" "}
							<ProfitLossBadge
								value={mf.day_profit_loss}
								percentage={mf.day_profit_loss_percentage}
							/>
						</span>
						<span className="px-2">
							Net:{" "}
							<ProfitLossBadge
								value={mf.net_profit_loss}
								percentage={mf.net_profit_loss_percentage}
							/>
						</span>
					</div>
				</div>
			</AccordionTrigger>

			<AccordionContent>
				<InvestmentTable mf={mf} />
			</AccordionContent>
		</AccordionItem>
	);
}
