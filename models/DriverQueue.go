package models

type DriverQueue struct {
	prev     *DriverQueue
	next     *DriverQueue
	RideInfo *Ride
}

func InitDriverQueue() *DriverQueue {
	newQueue := DriverQueue{}
	newQueue.next = &newQueue
	newQueue.prev = &newQueue
	return &newQueue
}

func (queue *DriverQueue) InsertRide(ride *Ride) {
	// If new list insert at the start, otherwise add a new node and insert the ride there
	if queue.RideInfo == nil {
		queue.RideInfo = ride
	} else {
		newQueue := DriverQueue{
			RideInfo: ride,
		}

		// If only 1 entry
		if queue.next == queue {
			newQueue.next = queue
			newQueue.prev = queue
			queue.next = &newQueue
			queue.prev = &newQueue
		} else {
			newQueue.prev = queue.prev
			queue.prev.next = &newQueue
			queue.prev = &newQueue
			newQueue.next = queue
		}
	}
}

func (queue *DriverQueue) GetRides() []Ride {
	var rides []Ride
	itQueue := queue

	for {
		if itQueue.RideInfo != nil {
			rides = append(rides, *itQueue.RideInfo)
		}

		itQueue = itQueue.next
		if itQueue == queue {
			break
		}
	}

	return rides
}

func (queue *DriverQueue) GetLastRide() *Ride {
	return queue.prev.RideInfo
}

func (queue *DriverQueue) UpdateLastRide(ride *Ride) {
	queue.prev.RideInfo = ride
}

func (queue *DriverQueue) PopRide() Ride {
	if queue.RideInfo != nil {
		toReturn := queue.RideInfo

		if queue.next == queue {
			queue.RideInfo = nil
			return *toReturn
		}

		place := queue

		*queue = *queue.next
		queue.next = place.next
		queue.prev = place.prev
		return *toReturn
	}
	return Ride{}
}

func (queue *DriverQueue) RemoveRide(ride *Ride) {
	i := 0
	itQueue := queue

	for {
		if itQueue.RideInfo.ID == ride.ID {
			if itQueue.next == queue && queue.Length() == 1 { // Only Item
				queue.RideInfo = nil
				break
			}

			if i == 0 {
				queue.PopRide()
				break
			}

			itQueue.prev.next = itQueue.next
			itQueue.next.prev = itQueue.prev
			itQueue.RideInfo = nil

			break
		}

		itQueue = itQueue.next
		if itQueue == queue {
			break
		}
		i++
	}
}

func (queue *DriverQueue) Length() int {
	i := 0
	itQueue := queue

	for {
		if itQueue.RideInfo != nil {
			i++
		}

		itQueue = itQueue.next
		if itQueue == queue {
			break
		}
	}

	return i
}
