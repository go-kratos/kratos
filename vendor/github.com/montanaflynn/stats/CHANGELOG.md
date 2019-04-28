# Change Log

## [0.2.0](https://github.com/montanaflynn/stats/tree/0.2.0)

### Merged pull requests:

- Fixed typographical error, changed accomdate to accommodate in README. [\#5](https://github.com/montanaflynn/stats/pull/5) ([saromanov](https://github.com/orthographic-pedant))

### Package changes:

- Add `Correlation` function
- Add `Covariance` function
- Add `StandardDeviation` function to be the same as `StandardDeviationPopulation`
- Change `Variance` function to be the same as `PopulationVariation`
- Add helper methods to `Float64Data`
- Add `Float64Data` type to use instead of `[]float64`
- Add `Series` type which references to `[]Coordinate`

## [0.1.0](https://github.com/montanaflynn/stats/tree/0.1.0)

Several functions were renamed in this release. They will still function but may be deprecated in the future.

### Package changes:

- Rename `VarP` to `PopulationVariance`
- Rename `VarS` to `SampleVariance`
- Rename `LinReg` to `LinearRegression`
- Rename `ExpReg` to `ExponentialRegression`
- Rename `LogReg` to `LogarithmicRegression`
- Rename `StdDevP` to `StandardDeviationPopulation`
- Rename `StdDevS` to `StandardDeviationSample`

## [0.0.9](https://github.com/montanaflynn/stats/tree/0.0.9)

### Closed issues:

- Functions have unexpected side effects [\#3](https://github.com/montanaflynn/stats/issues/3)
- Percentile is not calculated correctly [\#2](https://github.com/montanaflynn/stats/issues/2)

### Merged pull requests:

- Sample [\#4](https://github.com/montanaflynn/stats/pull/4) ([saromanov](https://github.com/saromanov))

### Package changes:

- Add HarmonicMean func
- Add GeometricMean func
- Add Outliers stuct and QuantileOutliers func
- Add Interquartile Range, Midhinge and Trimean examples
- Add Trimean
- Add Midhinge
- Add Inter Quartile Range
- Add Quantiles struct and Quantile func
- Add Nearest Rank method of calculating percentiles
- Add errors for all functions
- Add sample
- Add Linear, Exponential and Logarithmic Regression 
- Add sample and population variance and deviation 
- Add Percentile and Float64ToInt 
- Add Round 
- Add Standard deviation 
- Add Sum 
- Add Min and Ma- x 
- Add Mean, Median and Mode 
