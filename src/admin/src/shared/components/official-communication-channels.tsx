"use client";

import { Card, CardContent, CardHeader, Divider, Link, Stack, Typography } from "@mui/material"
import { SocialIcon } from 'react-social-icons/component'
import 'react-social-icons/github'
import 'react-social-icons/facebook'
import 'react-social-icons/bsky.app'
import 'react-social-icons/x'
import 'react-social-icons/linkedin'
import 'react-social-icons/discord'
import 'react-social-icons/reddit'
import 'react-social-icons/youtube'
import 'react-social-icons/tiktok'
import 'react-social-icons/threads'


export const OfficialCommunicationChannels = () => {

    return (
        <Card>
            <CardHeader title="Get help">
            </CardHeader>

            <CardContent>

                <Stack direction="row" spacing={0} gap={1} my={0}>

                    <Stack spacing={1}>


                        <Typography>
                            Read the
                            {' '}
                            <Link
                                href="#"
                                target="_blank"
                                variant="inherit"
                                color="action.disabled" 
								title="Not implemented!"
								sx={{ 
									cursor: 'not-allowed'
								}}
                            >
                                docs
                            </Link>

                        </Typography>


                        <Typography>
                            Check
                            {' '}
                            <Link
                                href="#"
                                target="_blank"
                                variant="inherit"
                                color="action.disabled"
								title="Not implemented!"
								sx={{ 
									cursor: 'not-allowed'
								}}
                            >
                                code examples
                            </Link>

                        </Typography>
                    </Stack>
                </Stack>
            </CardContent>

            <Divider />

            <CardContent>

                <Stack direction="row" flexWrap="wrap" spacing={0} gap={1} my={0}>
                    <SocialIcon url="https://github.com/GoLabra/labrago" target="_blank" style={{ height: 25, width: 25 }} rel="noopener"/>
                    {/* <SocialIcon url="https://linkedin.com/in/couetilc" target="_blank" style={{ height: 25, width: 25 }} rel="noopener" />
                    <SocialIcon url="https://facebook.com/" target="_blank" style={{ height: 25, width: 25 }} rel="noopener"/>
                    <SocialIcon url="https://bluesky.com/" network="bsky.app" style={{ height: 25, width: 25 }} rel="noopener"/>
                    <SocialIcon url="https://x.com/" style={{ height: 25, width: 25 }} rel="noopener"/>
                    <SocialIcon url="https://threads.com/" style={{ height: 25, width: 25 }} rel="noopener"/>
                    <SocialIcon url="https://discord.com/" target="_blank" style={{ height: 25, width: 25 }} rel="noopener"/>
                    <SocialIcon url="https://reddit.com/" target="_blank" style={{ height: 25, width: 25 }} rel="noopener"/>
                    <SocialIcon url="https://youtube.com/" target="_blank" style={{ height: 25, width: 25 }} rel="noopener"/>
                    <SocialIcon url="https://tiktok.com/" target="_blank" style={{ height: 25, width: 25 }} rel="noopener"/> */}

                </Stack>
            </CardContent>

            <Divider />


        </Card>
    )
}