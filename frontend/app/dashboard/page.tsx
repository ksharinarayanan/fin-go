"use client";

import React, { useEffect, useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { TrendingUp, TrendingDown, Calendar } from "lucide-react";
import { MutualFundInvestment, MutualFundInvestmentsResponse } from "./types";

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

const MutualFundCard = ({ mf }: { mf: MutualFundInvestment }) => {
	return (
		<Card className="overflow-hidden transition-shadow hover:shadow-lg">
			<CardHeader className="bg-secondary">
				<CardTitle className="text-lg font-bold">
					{mf.scheme_name}
				</CardTitle>
			</CardHeader>
			<CardContent className="p-4">
				<div className="grid grid-cols-2 gap-4">
					<div>
						<p className="text-sm text-muted-foreground">Units</p>
						<p className="text-lg font-semibold">
							{mf.units.toFixed(2)}
						</p>
					</div>
					<div>
						<p className="text-sm text-muted-foreground">
							Invested On
						</p>
						<p className="text-lg font-semibold flex items-center">
							<Calendar className="w-4 h-4 mr-1" />
							{new Date(mf.invested_at).toLocaleDateString()}
						</p>
					</div>
					<div>
						<p className="text-sm text-muted-foreground">
							Current NAV
						</p>
						<p className="text-lg font-semibold">
							{formatCurrency(mf.current_nav)}
						</p>
					</div>
					<div>
						<p className="text-sm text-muted-foreground">
							Invested NAV
						</p>
						<p className="text-lg font-semibold">
							{formatCurrency(mf.invested_nav)}
						</p>
					</div>
					<div className="col-span-2">
						<p className="text-sm text-muted-foreground">
							Current Value
						</p>
						<p className="text-xl font-bold flex items-center">
							{formatCurrency(mf.current_value)}
						</p>
					</div>
					<div className="col-span-2">
						<p className="text-sm text-muted-foreground">
							Invested Value
						</p>
						<p className="text-xl font-bold flex items-center">
							{formatCurrency(mf.invested_value)}
						</p>
					</div>
					<div className="col-span-2">
						<p className="text-sm text-muted-foreground">
							Day Profit/Loss
						</p>
						<p className="text-lg font-semibold flex items-center">
							<ProfitLossBadge
								value={mf.day_profit_loss}
								percentage={mf.day_profit_loss_percentage}
							/>
						</p>
					</div>
					<div className="col-span-2">
						<p className="text-sm text-muted-foreground">
							Net Profit/Loss
						</p>
						<p className="text-lg font-semibold flex items-center">
							<ProfitLossBadge
								value={mf.net_profit_loss}
								percentage={mf.net_profit_loss_percentage}
							/>
						</p>
					</div>
				</div>
			</CardContent>
		</Card>
	);
};

const Dashboard = () => {
	const [mfInvestments, setMfInvestments] = useState<MutualFundInvestment[]>(
		[]
	);
	const [isLoading, setIsLoading] = useState(true);
	const [error, setError] = useState<string | null>(null);

	useEffect(() => {
		const fetchData = async () => {
			try {
				const response = await fetch("http://localhost:8080/api/mf/");
				if (!response.ok) {
					throw new Error("Network response was not ok");
				}
				const data: MutualFundInvestmentsResponse =
					await response.json();
				setMfInvestments(data.mutual_funds);
			} catch (error) {
				setError("Failed to fetch mutual fund data");
				console.error("Error fetching mutual fund data:", error);
			} finally {
				setIsLoading(false);
			}
		};

		fetchData();
	}, []);

	if (isLoading) {
		return <div className="text-center py-10">Loading...</div>;
	}

	if (error) {
		return <div className="text-center py-10 text-red-500">{error}</div>;
	}

	return (
		<div className="container mx-auto p-4">
			<h1 className="text-3xl font-bold mb-6">Mutual Fund Investments</h1>
			<div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
				{mfInvestments.map((mf: MutualFundInvestment) => (
					<MutualFundCard key={mf.scheme_id} mf={mf} />
				))}
			</div>
		</div>
	);
};

export default Dashboard;
