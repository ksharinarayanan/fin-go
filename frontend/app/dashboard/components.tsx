import { Calendar, TrendingDown, TrendingUp } from "lucide-react";
import { MutualFundInvestment, MutualFundInvestmentDetails } from "./types";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import {
	AccordionContent,
	AccordionItem,
	AccordionTrigger,
} from "@/components/ui/accordion";
import { Table, TableCaption } from "@/components/ui/table";

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
			{formatCurrency(Math.abs(value))} ({percentage.toFixed(2)}%)
		</Badge>
	);
};

const MutualFundCard = ({
	investment,
}: {
	investment: MutualFundInvestmentDetails;
}) => {
	return (
		<Card className="overflow-hidden transition-shadow hover:shadow-lg">
			{/* <CardHeader className="bg-secondary">
				<CardTitle className="text-lg font-bold">
					{mf.scheme_name}
				</CardTitle>
			</CardHeader> */}
			<CardContent className="p-4">
				<div className="grid grid-cols-2 gap-4">
					<div>
						<p className="text-sm">Units</p>
						<p className="text-lg font-semibold">
							{investment.units.toFixed(2)}
						</p>
					</div>
					<div>
						<p className="text-sm">Invested On</p>
						<p className="text-lg font-semibold flex items-center">
							<Calendar className="w-4 h-4 mr-1" />
							{new Date(
								investment.invested_at
							).toLocaleDateString()}
						</p>
					</div>
					<div>
						<p className="text-sm">Invested NAV</p>
						<p className="text-lg font-semibold">
							{formatCurrency(investment.invested_nav)}
						</p>
					</div>
					<div className="col-span-2">
						<p className="text-sm">Current Value</p>
						<p className="text-xl font-bold flex items-center">
							{formatCurrency(investment.current_value)}
						</p>
					</div>
					<div className="col-span-2">
						<p className="text-sm">Invested Value</p>
						<p className="text-xl font-bold flex items-center">
							{formatCurrency(investment.invested_value)}
						</p>
					</div>
					<div className="col-span-2">
						<p className="text-sm">Net Profit/Loss</p>
						<div className="text-lg font-semibold flex items-center">
							<ProfitLossBadge
								value={investment.net_profit_loss}
								percentage={
									investment.net_profit_loss_percentage
								}
							/>
						</div>
					</div>
				</div>
			</CardContent>
		</Card>
	);
};

// please
function InvestmentTable() {}

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
				<Table>
					<TableCaption>{mf.scheme_name}</TableCaption>
				</Table>
				{mf.investments.map((investment, i) => {
					return <MutualFundCard key={i} investment={investment} />;
				})}
			</AccordionContent>
		</AccordionItem>
	);
}
