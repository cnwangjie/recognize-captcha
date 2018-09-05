package main

type point struct {
	data string
	label byte
}

type dataset struct {
	points []point
	labelsum int
}

func dist(a, b string) float64 {
	return (100.0 - float64(similarText(a, b))) / 100
}

func predict(dataset dataset, str string) byte {
	k := dataset.labelsum + 1
	labels := make([]byte, k)
	dists := make([]float64, k)
	for i := 0; i < k; i++ {
		dists[i] = 1
	}
	l := len(dataset.points)
	for i := 0; i < l; i++ {
		d := dist(dataset.points[i].data, str)
		if d > dists[k - 1] {
			continue
		}
		dists[k - 1] = d
		labels[k - 1] = dataset.points[i].label
		for j := k - 1; j > 0; j-- {
			if dists[j] > dists[j - 1] {
				break
			}
			dists[j], dists[j - 1] = dists[j - 1], dists[j]
			labels[j], labels[j - 1] = labels[j - 1], labels[j]
		}
	}
	// sum := make(map[byte]int)
	// mn := 0
	// var mb byte
	// for i := 0; i < k; i++ {
	// 	_, e := sum[labels[i]]
	// 	if e {
	// 		sum[labels[i]]++
	// 	} else {
	// 		sum[labels[i]] = 1
	// 	}
	// 	if sum[labels[i]] > mn {
	// 		mn = sum[labels[i]]
	// 		mb = labels[i]
	// 	}
	// }
	// return mb
	return labels[0]
}
