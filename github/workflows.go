package github

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
)

// ListWorkflows lists workflows in a repository
func (c *Client) ListWorkflows(ctx context.Context, owner, repo string, page, perPage int) ([]*Workflow, error) {
	url := fmt.Sprintf("repos/%s/%s/actions/workflows", owner, repo)
	params := []string{
		"page=" + strconv.Itoa(page),
		"per_page=" + strconv.Itoa(perPage),
	}
	url += "?" + params[0] + "&" + params[1]

	req, err := c.newRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response struct {
		TotalCount int         `json:"total_count"`
		Workflows  []*Workflow `json:"workflows"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return response.Workflows, nil
}

// ListWorkflowRuns lists workflow runs for a specific workflow
func (c *Client) ListWorkflowRuns(ctx context.Context, owner, repo string, workflowID int64, page, perPage int) ([]*WorkflowRun, error) {
	url := fmt.Sprintf("repos/%s/%s/actions/workflows/%d/runs", owner, repo, workflowID)
	params := []string{
		"page=" + strconv.Itoa(page),
		"per_page=" + strconv.Itoa(perPage),
	}
	url += "?" + params[0] + "&" + params[1]

	req, err := c.newRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response struct {
		TotalCount   int            `json:"total_count"`
		WorkflowRuns []*WorkflowRun `json:"workflow_runs"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return response.WorkflowRuns, nil
}

// TriggerWorkflow triggers a workflow run
func (c *Client) TriggerWorkflow(ctx context.Context, owner, repo string, workflowID int64, req *TriggerWorkflowRequest) error {
	url := fmt.Sprintf("repos/%s/%s/actions/workflows/%d/dispatches", owner, repo, workflowID)

	request, err := c.newRequest(ctx, "POST", url, req)
	if err != nil {
		return err
	}

	resp, err := c.do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// GetWorkflowRun gets a specific workflow run
func (c *Client) GetWorkflowRun(ctx context.Context, owner, repo string, runID int64) (*WorkflowRun, error) {
	url := fmt.Sprintf("repos/%s/%s/actions/runs/%d", owner, repo, runID)

	req, err := c.newRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var run WorkflowRun
	if err := json.NewDecoder(resp.Body).Decode(&run); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &run, nil
}

// CancelWorkflowRun cancels a workflow run
func (c *Client) CancelWorkflowRun(ctx context.Context, owner, repo string, runID int64) error {
	url := fmt.Sprintf("repos/%s/%s/actions/runs/%d/cancel", owner, repo, runID)

	req, err := c.newRequest(ctx, "POST", url, nil)
	if err != nil {
		return err
	}

	resp, err := c.do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// RerunWorkflow reruns a workflow
func (c *Client) RerunWorkflow(ctx context.Context, owner, repo string, runID int64) error {
	url := fmt.Sprintf("repos/%s/%s/actions/runs/%d/rerun", owner, repo, runID)

	req, err := c.newRequest(ctx, "POST", url, nil)
	if err != nil {
		return err
	}

	resp, err := c.do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
