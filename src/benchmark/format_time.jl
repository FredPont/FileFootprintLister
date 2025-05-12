using CSV, DataFrames
using Statistics
using VegaLite
path = "/mnt/crct-share/CRCT_P.Techno/Proteomique/Fred_P/2025/02_logiciels/logiciels_fred/FileFootprintLister/src/benchmark/benchmark_algos.csv"


function main()
	# Read the CSV file into a DataFrame
	df = CSV.read(path, DataFrame; delim = '\t', header = true)

	# Convert the time strings to seconds
	df.T1 = map(TM, df[:, 2])
	df.T2 = map(TM, df[:, 3])
	df.T3 = map(TM, df[:, 4])
	# Calculate the mean of columns A and B for each row
	#df.mean = mean.(eachrow(df[:, [:T1, :T2, :T4]]))

	# Print the DataFrame
	print(df)
	# Save the DataFrame to a new CSV file
	CSV.write("output.csv", df; delim = '\t', header = true)

	# Create a box plot
	box_plot = @vlplot(
		:boxplot,
		data = df,
		x = :category,
		y = :values,
		width = 400,
		height = 300
	)

	# Display the box plot
	display(box_plot)
end



# This function takes a string in the format "hh:mm:ss" and converts it to seconds
# It returns the total number of seconds as an integer
function TM(s)
	res = split(s, 'm')
	ti = parse(Float64, res[1]) * 60 + parse(Float64, res[2])
	return ti
end

main()