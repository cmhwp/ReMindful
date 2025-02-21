package algorithm

import (
	"math"
	"time"
)

// SM2Parameters SM-2算法的参数
type SM2Parameters struct {
	ReviewCount int       // 复习次数
	Difficulty  float64   // 难度系数 (0.8-5.0)
	LastReview  time.Time // 上次复习时间
	NextReview  time.Time // 下次复习时间
}

// Quality 复习质量评分
const (
	Complete    = 5 // 完全记住
	Correct     = 4 // 正确，有点困难
	Difficult   = 3 // 正确，但很困难
	Wrong       = 2 // 错误，但记得一些
	WrongHard   = 1 // 错误，有点印象
	WrongForget = 0 // 完全不记得
)

// CalculateNextReview 计算下次复习时间和更新难度系数
func CalculateNextReview(params *SM2Parameters, quality int) *SM2Parameters {
	// 更新难度系数 (0.8-5.0)
	diff := params.Difficulty + (0.1 - (5-float64(quality))*(0.08+(5-float64(quality))*0.02))
	diff = math.Max(0.8, math.Min(5.0, diff))

	// 计算间隔
	var interval time.Duration
	reviewCount := params.ReviewCount + 1

	switch {
	case quality < 3: // 答错了，重置复习计数
		interval = 24 * time.Hour
		reviewCount = 1
	case reviewCount == 1:
		interval = 24 * time.Hour
	case reviewCount == 2:
		interval = 6 * 24 * time.Hour
	default:
		// 间隔天数 = 复习次数 * 难度系数 * 基础间隔(1天)
		days := float64(reviewCount) * diff
		interval = time.Duration(days * float64(24*time.Hour))
	}

	now := time.Now()
	return &SM2Parameters{
		ReviewCount: reviewCount,
		Difficulty:  diff,
		LastReview:  now,
		NextReview:  now.Add(interval),
	}
}

// GetReviewQuality 根据用户表现返回复习质量
func GetReviewQuality(duration time.Duration, isCorrect bool, isHard bool) int {
	if !isCorrect {
		if duration < 5*time.Second {
			return WrongForget // 快速答错，说明完全不记得
		}
		if isHard {
			return WrongHard // 答错且觉得困难
		}
		return Wrong // 普通答错
	}

	if isHard {
		return Difficult // 答对但困难
	}
	if duration < 10*time.Second {
		return Complete // 快速答对，说明记得很清楚
	}
	return Correct // 普通答对
}
