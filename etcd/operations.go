package etcd

import (
	"context"
	"github.com/watcher-io/watcher/model"
	"go.etcd.io/etcd/clientv3"
)

func FetchMember(
	ctx context.Context,
	c *clientv3.Client,
) (
	*model.Cluster,
	error,
) {
	var clusterState model.Cluster
	memberListResponse, err := c.MemberList(ctx)
	if err != nil {
		return nil, err
	}

	clusterState.ID = memberListResponse.Header.ClusterId
	for _, member := range memberListResponse.Members {
		statusResponse, err := c.Status(ctx, member.ClientURLs[0])
		if err != nil {
			return nil, err
		}
		clusterState.Leader = statusResponse.Leader
		clusterState.Members = append(clusterState.Members, model.ClusterMember{
			ID:         member.ID,
			Name:       member.Name,
			PeerURLS:   member.PeerURLs,
			ClientURLS: member.ClientURLs,
			IsLearner:  member.IsLearner,
			Status: model.MemberStatus{
				Version:          statusResponse.Version,
				DbSize:           statusResponse.DbSize,
				DbSizeInUse:      statusResponse.DbSizeInUse,
				RaftIndex:        statusResponse.RaftIndex,
				RaftTerm:         statusResponse.RaftTerm,
				RaftAppliedIndex: statusResponse.RaftAppliedIndex,
			},
		})
	}
	return &clusterState, err
}

func PutKV(
	ctx context.Context,
	c *clientv3.Client,
	putKVRequest *model.PutKVRequest,
) (
	*model.PutKVResponse,
	error,
) {

	putResponse, err := c.Put(ctx, putKVRequest.Key, putKVRequest.Value)
	if err != nil {
		return nil, err
	}
	return &model.PutKVResponse{
		Revision: putResponse.Header.GetRevision(),
		MemberID: putResponse.Header.GetMemberId(),
		RaftTerm: putResponse.Header.GetRaftTerm(),
	}, nil
}
